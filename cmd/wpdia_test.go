package cmd

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	page = Page{
		Title:   "Golang",
		Extract: "Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson.",
	}
)

func TestWikiRequestBuilder(t *testing.T) {
	type args struct {
		params    url.Values
		baseURL   string
		userAgent string
	}
	tests := []struct {
		desc    string
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			desc: "Empty params arguments",
			args: args{
				params:    url.Values{},
				baseURL:   "https://api.example.com",
				userAgent: "Custom/User-Agent",
			},
			want: &http.Request{
				Method: "GET",
				Host:   "api.example.com",
				URL: &url.URL{
					Scheme:   "https",
					Host:     "api.example.com",
					RawQuery: "action=query&format=json",
				},
				Header: map[string][]string{
					"User-Agent":   {"Custom/User-Agent"},
					"Content-Type": {"multipart/form-data"},
				},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
			wantErr: false,
		},
		{
			desc: "Non empty params arguments",
			args: args{
				params: url.Values{
					"srlimit": {"1"},
					"utf8":    {"1"},
				},
				baseURL:   "https://api.example.com",
				userAgent: "Custom/User-Agent",
			},
			want: &http.Request{
				Method: "GET",
				Host:   "api.example.com",
				URL: &url.URL{
					Scheme:   "https",
					Host:     "api.example.com",
					RawQuery: "action=query&format=json&srlimit=1&utf8=1",
				},
				Header: map[string][]string{
					"User-Agent":   {"Custom/User-Agent"},
					"Content-Type": {"multipart/form-data"},
				},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
			wantErr: false,
		},
		{
			desc: "Empty User Agent",
			args: args{
				params:    url.Values{},
				baseURL:   "https://api.example.com",
				userAgent: "",
			},
			want: &http.Request{
				Method: "GET",
				Host:   "api.example.com",
				URL: &url.URL{
					Scheme:   "https",
					Host:     "api.example.com",
					RawQuery: "action=query&format=json",
				},
				Header: map[string][]string{
					"User-Agent":   {""},
					"Content-Type": {"multipart/form-data"},
				},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
			wantErr: false,
		},
		{
			desc: "Empty base URL",
			args: args{
				params:    url.Values{},
				baseURL:   "",
				userAgent: "Custom/User-Agent",
			},
			want: &http.Request{
				Method: "GET",
				Host:   "",
				URL: &url.URL{
					Scheme:   "",
					Host:     "",
					RawQuery: "action=query&format=json",
				},
				Header: map[string][]string{
					"User-Agent":   {"Custom/User-Agent"},
					"Content-Type": {"multipart/form-data"},
				},
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := wikiRequestBuilder(tt.args.params, tt.args.baseURL, tt.args.userAgent)

			// Copy the context from got to tt.want to avoid differences in deep equal assertion
			tt.want = tt.want.WithContext(got.Context())

			assert.EqualValues(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}
