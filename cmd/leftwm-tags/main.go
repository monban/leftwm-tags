package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/monban/leftwm-tags/internal/sexp"
)

type Root struct {
	WindowTitle string `json:"window_title"`
	Workspaces  []Workspace
}

type Workspace struct {
	X      uint   `json:"x"`
	Y      uint   `json:"y"`
	H      uint   `json:"h"`
	W      uint   `json:"w"`
	Output string `json:"output"`
	Index  uint   `json:"index"`
	Layout string `json:"layout"`
	Tags   []Tag
}

type Tag struct {
	Visible bool   `json:"visible"`
	Index   uint   `json:"index"`
	Busy    bool   `json:"busy"`
	Name    string `json:"name"`
	Mine    bool   `json:"mine"`
	Focused bool   `json:"focused"`
	Urgent  bool   `json:"urgent"`
}

func main() {
	data := bufio.NewScanner(os.Stdin)
	for data.Scan() {
		// Parse state JSON
		root := Root{}
		json.Unmarshal(data.Bytes(), &root)

		// Create widget
		box := sexp.NewList()
		box.Properties["class"] = sexp.String("workspace")
		box.Properties["orientation"] = sexp.String("h")
		box.Properties["space-evenly"] = sexp.Bool(true)
		box.Properties["halign"] = sexp.String("start")
		box.Properties["valign"] = sexp.String("fill")
		box.Properties["vexpand"] = sexp.Bool(true)
		box.Properties["spacing"] = sexp.Int(2)
		box.Elements = []sexp.Marshaler{sexp.Symbol("box")}
		for _, tag := range root.Workspaces[0].Tags {
			button := sexp.NewList()
			button.Elements = []sexp.Marshaler{sexp.Symbol("button"), sexp.String(tag.Name)}
			button.Properties["vexpand"] = sexp.Bool(true)
			button.Properties["hexpand"] = sexp.Bool(true)
			button.Properties["onclick"] = sexp.String(fmt.Sprintf("wmctrl -s %d", tag.Index))
			button.Properties["width"] = sexp.Int(32)

			// Set CSS classes
			classes := []string{}
			if tag.Focused {
				classes = append(classes, "active")
				button.Properties["active"] = sexp.Bool(false)
			}

			if tag.Busy {
				classes = append(classes, "busy")
			}

			if len(classes) > 0 {
				button.Properties["class"] = sexp.String(strings.Join(classes, " "))
			}
			box.Elements = append(box.Elements, button)
		}
		str, err := box.MarshalSexp()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Output widget
		os.Stdout.Write(str)
		os.Stdout.WriteString("\n")
	}
}
