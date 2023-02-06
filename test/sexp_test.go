package sexp

import (
	"testing"

	"github.com/matryer/is"
	"github.com/monban/leftwm-tags/internal/sexp"
)

func TestMarshallSymbol(t *testing.T) {
	i := is.New(t)
	expected := []byte("foo")
	in := sexp.Symbol("foo")
	result, err := in.MarshalSexp()
	i.NoErr(err)
	i.Equal(result, expected)
}

func TestMarshallString(t *testing.T) {
	i := is.New(t)
	expected := []byte("\"foo: \\\"bar\\\"\"")
	in := sexp.String("foo: \"bar\"")
	result, err := in.MarshalSexp()
	i.NoErr(err)
	i.Equal(string(result), string(expected))
}

func TestMarshalEmptyList(t *testing.T) {
	i := is.New(t)
	in := sexp.NewList()
	expected := "()"
	result, err := in.MarshalSexp()
	i.NoErr(err)
	i.Equal(string(result), expected)
}

func TestMarshallList(t *testing.T) {
	i := is.New(t)

	in := sexp.NewList()
	in.Elements = []sexp.Marshaler{sexp.Symbol("add"), sexp.Int(1), sexp.String("foo"), sexp.Bool(true)}
	in.Properties["str"] = sexp.String("bar")
	in.Properties["bool"] = sexp.Bool(false)
	in.Properties["integer"] = sexp.Int(42)

	expected := "(add :bool false :integer 42 :str \"bar\" 1 \"foo\" true)"

	result, err := in.MarshalSexp()

	i.NoErr(err)
	i.Equal(string(result), expected)
}
