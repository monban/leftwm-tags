package sexp

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Marshaler interface {
	MarshalSexp() ([]byte, error)
}

type Symbol string

func (s Symbol) MarshalSexp() ([]byte, error) {
	return []byte(s), nil
}

type String string

func (s String) MarshalSexp() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteRune('"')
	buf.WriteString(strings.ReplaceAll(string(s), "\"", "\\\""))
	buf.WriteRune('"')
	return buf.Bytes(), nil
}

type Int int

func (i Int) MarshalSexp() ([]byte, error) {
	return []byte(strconv.Itoa(int(i))), nil
}

type Bool bool

func (b Bool) MarshalSexp() ([]byte, error) {
	out := bytes.Buffer{}
	_, err := fmt.Fprintf(&out, "%t", b)
	return out.Bytes(), err
}

type List struct {
	Elements   []Marshaler
	Properties map[string]Marshaler
}

type Property struct {
	K string
	V Marshaler
}

func (p Property) WriteBuf(b *bytes.Buffer) error {
	value, err := p.V.MarshalSexp()
	if err != nil {
		return fmt.Errorf("Unable to marshal property of %s, value %s", p.K, p.V)
	}
	fmt.Fprintf(b, ":%s %s", p.K, value)
	return nil
}

func NewList() List {
	return List{
		Elements:   make([]Marshaler, 0),
		Properties: make(map[string]Marshaler, 0),
	}
}

func (l List) MarshalSexp() ([]byte, error) {
	if len(l.Elements) == 0 {
		return []byte("()"), nil
	}

	// Initialize the list
	buf := bytes.Buffer{}
	buf.WriteRune('(')

	// Push the first element into the output
	firstElm, err := l.Elements[0].MarshalSexp()
	if err != nil {
		return []byte{}, fmt.Errorf("Unable to marshal list: %v", err)
	}
	buf.Write(firstElm)

	// Write the properties
	properties := make([]string, len(l.Properties))
	i := 0
	for k := range l.Properties {
		properties[i] = k
		i++
	}
	sort.Strings(properties)

	for _, property := range properties {
		buf.WriteRune(' ')
		Property{property, l.Properties[property]}.WriteBuf(&buf)
	}

	// Collect remaining elements and separate by space
	ary := make([]string, len(l.Elements))
	for i := 1; i < len(l.Elements); i++ {
		j, err := l.Elements[i].MarshalSexp()
		if err != nil {
			return nil, err
		}
		ary[i] = string(j)
	}

	// Push the remaining elements
	buf.WriteString(strings.Join(ary, " "))

	// Close the list
	buf.WriteRune(')')
	return buf.Bytes(), nil
}
