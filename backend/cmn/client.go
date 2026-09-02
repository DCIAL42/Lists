package cmn

import (
	"context"
	"net/http"
)

type Client interface {
	BuildURL(map[string]string) string
	TryRequest(context.Context, string) (*http.Response, error)
	ReadToSearchResult(*http.Response, string) (SearchResult, error)
	Search(ctx context.Context, params map[string]string) (SearchResult, error)
	GetMedia(uint) (MediaResponse, error)
	ResolveMedia(Media) (MediaResponse, error)
}

type SearchResult struct {
	Next  string          `json:"next"`
	Items []MediaResponse `json:"items"`
}
