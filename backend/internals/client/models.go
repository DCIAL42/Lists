package client

import (
	"context"
	"net/http"
)

type Client struct {
	httpClient         *http.Client
	baseURL            string
	resultType         string
	configParams       map[string]string
	headers            map[string]string
	readToSearchResult func(*http.Response) (SearchResult, error)
	buildURL           func(string, map[string]string) string
	fetchToken         func(context.Context, *Client)
}

type SearchResult struct {
	Next  string `json:"next"`
	Items []any  `json:"items"`
}
