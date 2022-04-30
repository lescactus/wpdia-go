package cmd

import "time"

// WikiSearchResponse represents the http response of the Wikipedia Search API.
// Documentation is found here: https://www.mediawiki.org/wiki/API:Search
type WikiSearchResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		Sroffset int    `json:"sroffset"`
		Continue string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Search []struct {
			Ns        int       `json:"ns"`
			Title     string    `json:"title"`
			Pageid    uint64    `json:"pageid"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"search"`
	} `json:"query"`
}

// WikiTextExtractResponse represents the Wikipedia's TextExtracts API response
// Documentation is found here: https://www.mediawiki.org/wiki/Extension:TextExtracts#API
type WikiTextExtractResponse struct {
	Batchcomplete string `json:"batchcomplete"`
	Query         struct {
		Pages map[string]Page `json:"pages"`
	} `json:"query"`
}

// Page represents the page section of the Wikipedia's TextExtracts API response
// Documentation is found here: https://www.mediawiki.org/wiki/Extension:TextExtracts#API
type Page struct {
	// Use a pointer for the fields, so that the zero value of the JSON type
	// can be differentiated from the missing value

	Pageid *int `json:"pageid,omitempty" yaml:"pageid,omitempty"`
	Ns     *int `json:"ns,omitempty" yaml:"ns,omitempty"`

	Title   string `json:"title"`
	Extract string `json:"extract"`
}
