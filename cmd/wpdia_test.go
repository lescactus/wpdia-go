package cmd

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log = logrus.New()
}

var (
	ns        = 0
	pageid    = 25039021
	nsPtr     = &ns
	pageidPtr = &pageid

	page = Page{
		Title:   "Golang",
		Extract: "Go is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson.",

		Ns:     nsPtr,
		Pageid: pageidPtr,

		PageProps: &WikiPageProps{
			Disambiguation:    nil,
			WikiBaseShortDesc: "WikiBaseShortDesc",
			WikiBaseItem:      "WikiBaseItem",
		},
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

func TestNewWikiClient(t *testing.T) {
	type args struct {
		baseURL   string
		userAgent string
	}
	tests := []struct {
		desc    string
		args    args
		want    *WikiClient
		wantErr bool
	}{
		{
			desc: "BaseURL and UserAgent are set and valid",
			args: args{baseURL: "https://api.example.com", userAgent: "Custom/User-Agent"},
			want: &WikiClient{
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "api.example.com",
				},
				UserAgent: "Custom/User-Agent" + defaultUserAgent,
				Client: &http.Client{
					Transport:     nil,
					CheckRedirect: nil,
					Jar:           nil,
					Timeout:       timeout,
				},
			},
			wantErr: false,
		},
		{
			desc: "BaseURL is set and valid, UserAgent is not set",
			args: args{baseURL: "https://api.example.com", userAgent: ""},
			want: &WikiClient{
				BaseURL: &url.URL{
					Scheme: "https",
					Host:   "api.example.com",
				},
				UserAgent: defaultUserAgent,
				Client: &http.Client{
					Transport:     nil,
					CheckRedirect: nil,
					Jar:           nil,
					Timeout:       timeout,
				},
			},
			wantErr: false,
		},
		{
			desc:    "BaseURL is set and invalid, UserAgent is not set",
			args:    args{baseURL: "\ninvalid url", userAgent: ""},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := NewWikiClient(tt.args.baseURL, tt.args.userAgent)

			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
