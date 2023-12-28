package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
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
	log.WithFields(logrus.Fields{
		"level": logLevel,
		"url":   baseURL,
	}).Debug("Parsing base URL...")

	// Ensure the base URL is valid
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"level": logLevel,
		"url":   baseURL,
	}).Debug("Base URL Parsed")

	// If empty, the User-Agent provided as a parameter will be set with a default value
	// If not empty, the User-Agent provided as a parameter will be merged to the default one
	var ua string
	if userAgent == "" {
		ua = defaultUserAgent
	} else {
		ua = userAgent + defaultUserAgent
	}

	log.WithFields(logrus.Fields{
		"level":      logLevel,
		"user-agent": ua,
	}).Debug("User-Agent set")

	return &WikiClient{
		BaseURL:   url,
		UserAgent: ua,
		Client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       timeout,
		},
	}, nil
}

// GetExtract will invoke the Wikipedia's TextExtracts's API to extract the text of the given page id.
// It takes in argument the page id to request and will return the response or any error encountered.
func (w *WikiClient) GetExtract(id uint64) (*WikiTextExtractResponse, error) {
	log.WithFields(logrus.Fields{
		"level": logLevel,
	}).Debug("Setting http request parameters...")

	params := wikiExtractRequestParamsBuilder(exintro)
	params.Add("pageids", fmt.Sprintf("%d", id))

	log.WithFields(logrus.Fields{
		"level":  logLevel,
		"params": params,
	}).Debug("Http request parameters set")

	return w.do(params)
}

// GetExtractRandom will invoke the Wikipedia's Random API to fetch the content of a random article.
// It takes no argument and will return the response or any error encountered.
func (w *WikiClient) GetExtractRandom() (*WikiTextExtractResponse, error) {
	log.WithFields(logrus.Fields{
		"level": logLevel,
	}).Debug("Setting http request parameters...")

	params := wikiExtractRequestParamsBuilder(exintro)

	// When requesting a random page,
	// 'genarator=random' parameter must be set
	// 'genarator=random' provide a set of random pages
	params.Add("generator", "random")
	// Namespace 0 is 'Articles'. ref: https://www.mediawiki.org/wiki/Manual:Namespace
	params.Add("grnnamespace", "0")
	// Limit to only 1 random page returned
	params.Add("grnlimit", "1")

	log.WithFields(logrus.Fields{
		"level":  logLevel,
		"params": params,
	}).Debug("Http request parameters set")

	return w.do(params)
}

// do will build a http request with the given http request parameters as arguments,
// execute it and unmarshal the response to a *WikiTextExtractResponse.
// It will use the embedded BaseURL and User-Agent.
// It will take care of reading the body response and to close it.
//
// The function takes as argument a set of url query parameters and will return the response or any error encountered.
func (w *WikiClient) do(params url.Values) (*WikiTextExtractResponse, error) {
	log.WithFields(logrus.Fields{
		"level":      logLevel,
		"params":     params,
		"url":        w.BaseURL.String(),
		"user-agent": w.UserAgent,
	}).Debug("Building http request...")

	// Build http request
	req, err := wikiRequestBuilder(params, w.BaseURL.String(), w.UserAgent)
	if err != nil {
		return nil, fmt.Errorf("error while building http request: %v", err)
	}

	log.WithFields(logrus.Fields{
		"level":      logLevel,
		"params":     params,
		"url":        w.BaseURL.String(),
		"user-agent": w.UserAgent,
	}).Debug("Http request built")

	log.WithFields(logrus.Fields{
		"level": logLevel,
	}).Debug("Sending http request...")

	// Execute the http request
	resp, err := w.Client.Do(req)
	if err != nil {
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"level": logLevel,
	}).Debug("Http request sent")

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.WithFields(logrus.Fields{
		"level": logLevel,
	}).Debug("Reading http response body and unmarshalling...")

	var r WikiTextExtractResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	log.WithFields(logrus.Fields{
		"level": logLevel,
	}).Debug("Http response body read and unmarshalled")

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
	params.Add("list", "search")
	params.Add("utf8", "1")
	params.Add("srsearch", title)

	log.WithFields(logrus.Fields{
		"level":  logLevel,
		"params": params,
	}).Debug("Http request parameters set")

	log.WithFields(logrus.Fields{
		"level":      logLevel,
		"params":     params,
		"url":        w.BaseURL.String(),
		"user-agent": w.UserAgent,
	}).Debug("Building http request...")

	// Build http request
	req, err := wikiRequestBuilder(params, w.BaseURL.String(), w.UserAgent)
	if err != nil {
		return 0, fmt.Errorf("error while building http request: %v", err)
	}

	log.WithFields(logrus.Fields{
		"level":      logLevel,
		"params":     params,
		"url":        w.BaseURL.String(),
		"user-agent": w.UserAgent,
	}).Debug("Http request built")

	log.WithFields(logrus.Fields{
		"level": logLevel,
	}).Debug("Sending http request...")

	resp, err := w.Client.Do(req)
	if err != nil {
		return 0, err
	}

	log.WithFields(logrus.Fields{
		"level": logLevel,
	}).Debug("Http request sent")

	log.WithFields(logrus.Fields{
		"level": logLevel,
	}).Debug("Reading http response body and unmarshalling...")

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var s WikiSearchResponse
	err = json.Unmarshal(body, &s)
	if err != nil {
		return 0, err
	}

	log.WithFields(logrus.Fields{
		"level": logLevel,
	}).Debug("Http response body read and unmarshalled")

	// Query.Search[] will be empty if the search doesn't match anything
	if len(s.Query.Search) == 0 {
		log.WithFields(logrus.Fields{
			"level": logLevel,
		}).Warn("Search didn't match anything")
		return 0, nil
	}

	log.WithFields(logrus.Fields{
		"level":  logLevel,
		"pageid": s.Query.Search[0].Pageid,
	}).Info("Search found a Page ID")

	// We only care about the first result
	return s.Query.Search[0].Pageid, nil
}

// wikiRequestBuilder is used to build a http request to the Wikipedia's API.
// It will create a http GET request with:
// - a set of standard http parameters in addition to the one passed to the function,
// - a User-Agent http header to follow the best practice and etiquette for the use of Wikipedia's API,
// - a valid Content-Type http header
//
// The function takes as argument a set of url query parameters, the base URL and the User-Agent.
// It returns a *http.Request or any error encountered.
func wikiRequestBuilder(params url.Values, baseURL, userAgent string) (*http.Request, error) {
	// Common parameters for each requests to Wikipedia API
	params.Add("action", "query")
	params.Add("format", "json")

	// URL encode the parameters
	req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	// According to https://www.mediawiki.org/wiki/API:Etiquette#The_User-Agent_header
	// "
	// It is best practice to set a descriptive User Agent header.
	// To do so, use User-Agent: clientname/version (contact information e.g. username, email) framework/version....
	// "
	req.Header.Set("User-Agent", userAgent)

	// According to https://www.mediawiki.org/wiki/API:Data_formats
	// "
	// The API takes its input through parameters provided by the HTTP request in
	// application/x-www-form-urlencoded or multipart/form-data format.
	// "
	req.Header.Set("Content-Type", "multipart/form-data")

	return req, nil
}

// wikiExtractRequestParamsBuilder is used to provide some base http parameters for the TextExtract API.
// It will create a url.Values with the required properties for extracting text.
// The documentation can be found here: https://www.mediawiki.org/wiki/Extension:TextExtracts#API
//
// The function takes as argument a boolean value indicating whether or not only requesting the content before the first section
// and will returns a url.Values.
func wikiExtractRequestParamsBuilder(exintro bool) url.Values {
	params := url.Values{}

	params.Add("explaintext", "1")
	params.Add("exsectionformat", "plain")
	params.Add("prop", "extracts|pageprops")

	// 'exintro' is mutually exclusive with 'exsentences'
	// Either we return only the content before the first section
	// or we return a given number of sentences
	if exintro {
		params.Add("exintro", "1")
	} else {
		params.Add("exsentences", exsentences)
	}

	return params
}
