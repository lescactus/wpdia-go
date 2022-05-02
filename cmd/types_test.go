package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	str = ""
)

func TestPageIsDisambiguation(t *testing.T) {

	tests := []struct {
		name string
		p    Page
		want bool
	}{
		{
			name: "Disambiguation is nil",
			p: Page{
				PageProps: &WikiPageProps{
					Disambiguation:    nil,
					WikiBaseShortDesc: "WikiBaseShortDesc",
					WikiBaseItem:      "",
				},
			},
			want: false,
		},
		{
			name: "Disambiguation is not nil",
			p: Page{
				PageProps: &WikiPageProps{
					Disambiguation:    &str,
					WikiBaseShortDesc: "WikiBaseShortDesc",
					WikiBaseItem:      "",
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.IsDisambiguation()

			assert.Equal(t, tt.want, got)
		})
	}
}
