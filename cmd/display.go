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
	Write(w io.Writer, p *Page, full bool) error
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

func (d *plainFormat) Write(w io.Writer, p *Page, full bool) error {
	_, err := fmt.Fprintf(w, "Title:\n  %s\n\n", p.Title)
	if err != nil {
		return err
	}

	if full {
		_, err := fmt.Fprintf(w, "Ns:\n  %d\n\n", *p.Ns)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "Pageid:\n  %d\n\n", *p.Pageid)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "WikiBase Short Description:\n  %s\n\n", p.PageProps.WikiBaseShortDesc)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "WikiBase Item:\n  %s\n\n", p.PageProps.WikiBaseItem)
		if err != nil {
			return err
		}

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

func (d *prettyFormat) Write(w io.Writer, p *Page, full bool) error {
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

	if full {
		out, err := r.Render("### Namespace" + fmt.Sprintf("\nNs: %d", *p.Ns))
		if err != nil {
			return err
		}
		fmt.Fprint(w, out)

		out, err = r.Render("### Page ID " + fmt.Sprintf("\nPageid: %d", *p.Pageid))
		if err != nil {
			return err
		}
		fmt.Fprint(w, out)

		out, err = r.Render("### WikiBase Short Description" + fmt.Sprintf("\nNs: %s", p.PageProps.WikiBaseShortDesc))
		if err != nil {
			return err
		}
		fmt.Fprint(w, out)

		out, err = r.Render("### WikiBase Item" + fmt.Sprintf("\nNs: %s", p.PageProps.WikiBaseItem))
		if err != nil {
			return err
		}
		fmt.Fprint(w, out)
	}

	out, err = r.Render(fmt.Sprintf("### Extract\n%s", p.Extract))
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

func (d *jsonFormat) Write(w io.Writer, p *Page, full bool) error {
	// Nullify these fields if not requesting the full output
	if !full {
		var n *int
		var i *int
		p.Ns = n
		p.Pageid = i

		p.PageProps = nil
	}

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

func (d *yamlFormat) Write(w io.Writer, p *Page, full bool) error {
	// Nullify these fields if not requesting the full output
	if !full {
		var n *int
		var i *int
		p.Ns = n
		p.Pageid = i

		p.PageProps = nil
	}

	out, err := yaml.Marshal(&p)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", string(out))

	return nil
}
