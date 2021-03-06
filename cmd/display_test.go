package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlainFormat(t *testing.T) {
	tests := []struct {
		desc string
		want *plainFormat
	}{
		{
			desc: "NewPlainFormat",
			want: &plainFormat{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.want, NewPlainFormat())
		})
	}
}

func TestNewPrettyFormat(t *testing.T) {
	type args struct {
		wordWrap int
	}
	tests := []struct {
		desc string
		args args
		want *prettyFormat
	}{
		{
			desc: "wordWrap = 100",
			args: args{wordWrap: 100},
			want: &prettyFormat{wordWrap: 100},
		},
		{
			desc: "wordWrap = 0",
			args: args{wordWrap: 0},
			want: &prettyFormat{wordWrap: 100},
		},
		{
			desc: "wordWrap = -1",
			args: args{wordWrap: -1},
			want: &prettyFormat{wordWrap: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.want, NewPrettyFormat(tt.args.wordWrap))
		})
	}
}

func TestNewJsonFormat(t *testing.T) {
	type args struct {
		prefix string
		indent string
	}
	tests := []struct {
		desc string
		args args
		want *jsonFormat
	}{
		{
			desc: "Prefix: '', Indent: '    '",
			args: args{prefix: "", indent: "    "},
			want: &jsonFormat{prefix: "", indent: "    "},
		},
		{
			desc: "Prefix: ' ', Indent: ''",
			args: args{prefix: " ", indent: ""},
			want: &jsonFormat{prefix: " ", indent: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.want, NewJsonFormat(tt.args.prefix, tt.args.indent))
		})
	}
}

func TestNewYamlFormat(t *testing.T) {
	tests := []struct {
		desc string
		want *yamlFormat
	}{
		{
			desc: "NewYamlFormat",
			want: &yamlFormat{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.want, NewYamlFormat())
		})
	}
}

func TestPlainFormatWrite(t *testing.T) {
	type args struct {
		p    *Page
		full bool
	}
	tests := []struct {
		name    string
		d       *plainFormat
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:    "Without full output",
			d:       NewPlainFormat(),
			args:    args{p: &page, full: false},
			wantW:   fmt.Sprintf("Title:\n  %s\n\nExtract:\n  %s", page.Title, page.Extract),
			wantErr: false,
		},
		{
			name: "With full output",
			d:    NewPlainFormat(),
			args: args{p: &page, full: true},
			wantW: fmt.Sprintf("Title:\n  %s\n\nNs:\n  %d\n\nPageid:\n  %d\n\nWikiBase Short Description:\n  %s\n\nWikiBase Item:\n  %s\n\nExtract:\n  %s",
				page.Title,
				*page.Ns,
				*page.Pageid,
				page.PageProps.WikiBaseShortDesc,
				page.PageProps.WikiBaseItem,
				page.Extract),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &plainFormat{}
			w := &bytes.Buffer{}
			if err := d.Write(w, tt.args.p, tt.args.full); (err != nil) != tt.wantErr {
				t.Errorf("plainFormat.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("plainFormat.Write() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPrettyFormatWrite(t *testing.T) {
	type fields struct {
		wordWrap int
	}
	type args struct {
		p    *Page
		full bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantW   string
		wantErr bool
	}{
		{
			fields:  fields{wordWrap: 100},
			args:    args{p: &page, full: true},
			wantErr: false,
		},
		{
			fields:  fields{wordWrap: 100},
			args:    args{p: &page, full: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			d := &prettyFormat{
				wordWrap: tt.fields.wordWrap,
			}
			w := &bytes.Buffer{}
			err := d.Write(w, tt.args.p, tt.args.full)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestJsonFormatWrite(t *testing.T) {
	type fields struct {
		prefix string
		indent string
	}
	type args struct {
		p    *Page
		full bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name:   "prefix: '', indent: '    ', full: false",
			fields: fields{prefix: "", indent: "    "},
			args:   args{full: false},
			wantW: fmt.Sprintf(`{
    "title": "%s",
    "extract": "%s"
}
`, page.Title, page.Extract),
			wantErr: false,
		},
		{
			name:   "prefix: ' ', indent: '', full: false",
			fields: fields{prefix: " ", indent: ""},
			args:   args{full: false},
			wantW: fmt.Sprintf(`{
 "title": "%s",
 "extract": "%s"
 }
`, page.Title, page.Extract),
			wantErr: false,
		},
		{
			name:   "prefix: '', indent: '    ', full: true",
			fields: fields{prefix: "", indent: "    "},
			args:   args{full: true},
			wantW: fmt.Sprintf(`{
    "pageid": %d,
    "ns": %d,
    "title": "%s",
    "extract": "%s",
    "pageprops": {
        "wikibase-shortdesc": "%s",
        "wikibase_item": "%s"
    }
}
`, *page.Pageid, *page.Ns, page.Title, page.Extract, page.PageProps.WikiBaseShortDesc, page.PageProps.WikiBaseItem),
			wantErr: false,
		},
		{
			name:   "prefix: ' ', indent: '', full: true",
			fields: fields{prefix: " ", indent: ""},
			args:   args{full: true},
			wantW: fmt.Sprintf(`{
 "pageid": %d,
 "ns": %d,
 "title": "%s",
 "extract": "%s",
 "pageprops": {
 "wikibase-shortdesc": "%s",
 "wikibase_item": "%s"
 }
 }
`, *page.Pageid, *page.Ns, page.Title, page.Extract, page.PageProps.WikiBaseShortDesc, page.PageProps.WikiBaseItem),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &jsonFormat{
				prefix: tt.fields.prefix,
				indent: tt.fields.indent,
			}
			w := &bytes.Buffer{}

			newPage := page
			tt.args.p = &newPage
			err := d.Write(w, tt.args.p, tt.args.full)

			assert.Equal(t, tt.wantW, w.String())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestYamlFormatWrite(t *testing.T) {
	type args struct {
		p    *Page
		full bool
	}
	tests := []struct {
		name    string
		d       *yamlFormat
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "Without full output",
			d:    &yamlFormat{},
			args: args{full: false},
			wantW: fmt.Sprintf(`title: %s
extract: Go is a statically typed, compiled programming language designed at Google
  by Robert Griesemer, Rob Pike, and Ken Thompson.

`, page.Title),
			wantErr: false,
		},
		{
			name: "With full output",
			d:    &yamlFormat{},
			args: args{full: true},
			wantW: fmt.Sprintf(`pageid: %d
ns: %d
title: %s
extract: Go is a statically typed, compiled programming language designed at Google
  by Robert Griesemer, Rob Pike, and Ken Thompson.
pageprops:
  wikibase-shortdesc: %s
  wikibase_item: %s

`, *page.Pageid, *page.Ns, page.Title, page.PageProps.WikiBaseShortDesc, page.PageProps.WikiBaseItem),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}

			newPage := page
			tt.args.p = &newPage
			err := tt.d.Write(w, tt.args.p, tt.args.full)

			assert.Equal(t, tt.wantW, w.String())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
