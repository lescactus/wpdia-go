package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasOneArg(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "One arg",
			args: args{[]string{"one"}},
			want: true,
		},
		{
			name: "Two args",
			args: args{[]string{"one", "two"}},
			want: false,
		},
		{
			name: "Three args",
			args: args{[]string{"one", "two", "three"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasOneArg(tt.args.args)

			assert.Equal(t, tt.want, got)
		})
	}
}
