package client

import (
	"context"
	"net/http"
)

type Client interface {
	BuildURL(map[string]string) string
	TryRequest(context.Context, string) (*http.Response, error)
	ReadToSearchResult(*http.Response) (SearchResult, error)
	Search(context.Context, map[string]string) (SearchResult, error)
	GetResultType() string
}

type SearchResult struct {
	Next  string `json:"next"`
	Items []any  `json:"items"`
}
