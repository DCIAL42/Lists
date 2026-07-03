package client

import (
	"net/http"
)

type Client struct {
	httpClient         *http.Client
	baseURL            string
	resultType         string
	configParams       map[string]string
	headers            map[string]string
	readToSearchResult func(*http.Response) ([]SearchResult, error)
}

type SearchResult any
