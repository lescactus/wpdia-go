package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/charmbracelet/glamour"
	"gopkg.in/yaml.v2"
)

// Displayer offers function to display a page
// using different formatters.
type Displayer interface {
	// Write will write the content of a page
	// to the given io.Writer
	Write(w io.Writer, p *Page) error
}

type plainFormat struct{}

type prettyFormat struct {
	wordWrap int
}

type jsonFormat struct {
	prefix string
	indent string
}

type yamlFormat struct{}

func NewPlainFormat() *plainFormat {
	return &plainFormat{}
}

func (d *plainFormat) Write(w io.Writer, p *Page) error {
	_, err := fmt.Fprintf(w, "Title:\n  %s\n\n", p.Title)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "Extract:\n  %s", p.Extract)
	if err != nil {
		return err
	}

	return nil
}

func NewPrettyFormat(wordWrap int) *prettyFormat {
	// set wordWrap to 100 by default
	if wordWrap <= 0 {
		wordWrap = 100
	}

	return &prettyFormat{wordWrap: wordWrap}
}

func (d *prettyFormat) Write(w io.Writer, p *Page) error {
	r, err := glamour.NewTermRenderer(
		// detect background color and pick either the default dark or light theme
		glamour.WithAutoStyle(),
		// wrap output at specific width
		glamour.WithWordWrap(d.wordWrap),
	)
	if err != nil {
		return err
	}

	out, err := r.Render("## " + p.Title)
	if err != nil {
		return err
	}
	fmt.Fprint(w, out)

	out, err = r.Render(p.Extract)
	if err != nil {
		return err
	}
	fmt.Fprint(w, out)

	return nil
}

func NewJsonFormat(prefix, indent string) *jsonFormat {
	return &jsonFormat{
		prefix: prefix,
		indent: indent,
	}
}

func (d *jsonFormat) Write(w io.Writer, p *Page) error {
	// Nullify these fields as we are not interested in Formating them
	p.Ns = 0
	p.Pageid = 0

	b, err := json.MarshalIndent(p, d.prefix, d.indent)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, string(b))

	return nil
}

func NewYamlFormat() *yamlFormat {
	return &yamlFormat{}
}

func (d *yamlFormat) Write(w io.Writer, p *Page) error {
	// Nullify these fields as we are not interested in Formating them
	p.Ns = 0
	p.Pageid = 0

	out, err := yaml.Marshal(&p)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", string(out))

	return nil
}
