package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	// defaultUserAgent is the http User-Agent used by default.
	//
	// The API etiquette of the MediaWiki API ask clients to provide an informative User-Agent.
	// The generic format is <client name>/<version> (<contact information>) <library/framework name>/<version> [<library name>/<version> ...]
	//
	// Ref: https://meta.wikimedia.org/wiki/User-Agent_policy
	defaultUserAgent = "wpdia-go/" + version + " (github.com/lescactus/wpdia-go) WikiClient/" + version
)

// WikiClient represents the API client
type WikiClient struct {
	BaseURL   *url.URL
	UserAgent string
	Client    *http.Client
}

// NewWikiClient creates a new WikiClient with a given API base URL and http User-Agent.
// When the User-Agent is empty, it uses a default one.
// It returns a WikiClient or any error encountered
func NewWikiClient(baseURL, userAgent string) (*WikiClient, error) {
	// Ensure the base URL is valid
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	// If empty, the User-Agent provided as a parameter will be set with a default value
	// If not empty, the User-Agent provided as a parameter will be merged to the default one
	var ua string
	if userAgent == "" {
		ua = defaultUserAgent
	} else {
		ua = userAgent + defaultUserAgent
	}

	return &WikiClient{
		BaseURL:   url,
		UserAgent: ua,
		Client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       15 * time.Second,
		},
	}, nil
}

func (w *WikiClient) GetExtract(id uint64) (*WikiTextExtractResponse, error) {
	params := url.Values{}

	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("prop", "extracts")
	params.Add("pageids", fmt.Sprint(id))
	params.Add("exsentences", "10")
	params.Add("explaintext", "1")
	params.Add("exsectionformat", "plain")

	req, err := http.NewRequest("GET", w.BaseURL.String()+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	// According to https://www.mediawiki.org/wiki/API:Etiquette#The_User-Agent_header
	// "
	// It is best practice to set a descriptive User Agent header.
	// To do so, use User-Agent: clientname/version (contact information e.g. username, email) framework/version....
	// "
	req.Header.Set("User-Agent", w.UserAgent)

	// According to https://www.mediawiki.org/wiki/API:Data_formats
	// "
	// The API takes its input through parameters provided by the HTTP request in
	// application/x-www-form-urlencoded or multipart/form-data format.
	// "
	req.Header.Set("Content-Type", "multipart/form-data")

	resp, err := w.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var r WikiTextExtractResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// SearchTitle will invoke the Wikipedia's Search API to lookup for the given title.
// It takes in argument the title to search for and will return the page id of the first
// result if found. If the search doesn't return any result, the function return 0 or
// any error encountered.
func (w *WikiClient) SearchTitle(title string) (uint64, error) {

	params := url.Values{}

	// Documentation about the search API: https://www.mediawiki.org/wiki/API:Search
	//
	// "srsearch" will search for page titles or page content
	// matching the given value.
	//
	// We only care about the first result of the search
	// which should match what we are searching for
	params.Add("srlimit", "1")
	params.Add("action", "query")
	params.Add("format", "json")
	params.Add("list", "search")
	params.Add("utf8", "1")
	params.Add("srsearch", title)

	req, err := http.NewRequest("GET", w.BaseURL.String()+"?"+params.Encode(), nil)
	if err != nil {
		return 0, err
	}

	// According to https://www.mediawiki.org/wiki/API:Etiquette#The_User-Agent_header
	// "
	// It is best practice to set a descriptive User Agent header.
	// To do so, use User-Agent: clientname/version (contact information e.g. username, email) framework/version....
	// "
	req.Header.Set("User-Agent", w.UserAgent)

	// According to https://www.mediawiki.org/wiki/API:Data_formats
	// "
	// The API takes its input through parameters provided by the HTTP request in
	// application/x-www-form-urlencoded or multipart/form-data format.
	// "
	req.Header.Set("Content-Type", "multipart/form-data")

	resp, err := w.Client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var s WikiSearchResponse
	err = json.Unmarshal(body, &s)
	if err != nil {
		return 0, err
	}

	// Query.Search[] will be empty if the search doesn't match anything
	if len(s.Query.Search) == 0 {
		return 0, nil
	}

	// We only care about the first result
	return s.Query.Search[0].Pageid, nil
}

func displayExtract(p Page) {
	fmt.Println(p.Title)
	fmt.Println(p.Extract)
}
