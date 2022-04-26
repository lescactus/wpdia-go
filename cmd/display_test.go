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
		p *Page
	}
	tests := []struct {
		desc    string
		d       *plainFormat
		args    args
		wantW   string
		wantErr bool
	}{
		{
			d:       NewPlainFormat(),
			args:    args{p: &page},
			wantW:   fmt.Sprintf("Title:\n  %s\n\nExtract:\n  %s", page.Title, page.Extract),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			d := &plainFormat{}
			w := &bytes.Buffer{}
			if err := d.Write(w, tt.args.p); (err != nil) != tt.wantErr {
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
		p *Page
	}
	tests := []struct {
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields:  fields{wordWrap: 100},
			args:    args{p: &page},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			d := &prettyFormat{
				wordWrap: tt.fields.wordWrap,
			}
			w := &bytes.Buffer{}
			err := d.Write(w, tt.args.p)

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
		p *Page
	}
	tests := []struct {
		desc    string
		fields  fields
		args    args
		wantW   string
		wantErr bool
	}{
		{
			desc:   "prefix: '', indent: '    '",
			fields: fields{prefix: "", indent: "    "},
			args:   args{p: &page},
			wantW: fmt.Sprintf(`{
    "title": "%s",
    "extract": "%s"
}
`, page.Title, page.Extract),
			wantErr: false,
		},
		{
			desc:   "prefix: ' ', indent: ''",
			fields: fields{prefix: " ", indent: ""},
			args:   args{p: &page},
			wantW: fmt.Sprintf(`{
 "title": "%s",
 "extract": "%s"
 }
`, page.Title, page.Extract),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			d := &jsonFormat{
				prefix: tt.fields.prefix,
				indent: tt.fields.indent,
			}
			w := &bytes.Buffer{}
			err := d.Write(w, tt.args.p)

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
		p *Page
	}
	tests := []struct {
		desc    string
		d       *yamlFormat
		args    args
		wantW   string
		wantErr bool
	}{
		{
			desc: "",
			d:    &yamlFormat{},
			args: args{p: &page},
			wantW: fmt.Sprintf(`title: %s
extract: Go is a statically typed, compiled programming language designed at Google
  by Robert Griesemer, Rob Pike, and Ken Thompson.

`, page.Title),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := tt.d.Write(w, tt.args.p)

			assert.Equal(t, tt.wantW, w.String())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
