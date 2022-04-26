package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	type args struct {
		s   []string
		str string
	}
	tests := []struct {
		desc string
		args args
		want bool
	}{
		{
			desc: "String is present in slice",
			args: args{
				s: []string{"foo", "bar"},
				str: "foo",
			},
			want: true,
		},
		{
			desc: "String is not present in slice",
			args: args{
				s: []string{"foo", "bar"},
				str: "xxx",
			},
			want: false,
		},
		{
			desc: "Slice is empty",
			args: args{
				s: []string{},
				str: "xxx",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// if got := isPresent(tt.args.s, tt.args.str); got != tt.want {
			// 	t.Errorf("isPresent() = %v, want %v", got, tt.want)
			// }
			res := isPresent(tt.args.s, tt.args.str)
			assert.Equal(t, tt.want, res)
		})
	}
	
}
