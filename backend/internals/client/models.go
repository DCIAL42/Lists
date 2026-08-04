package client

import (
	"context"
	"net/http"

	"github.com/DCIAL42/media/cmn"
)

type Client interface {
	BuildURL(map[string]string) string
	TryRequest(context.Context, string) (*http.Response, error)
	ReadToSearchResult(*http.Response) (SearchResult, error)
	Search(context.Context, map[string]string) (SearchResult, error)
	GetMediaType() cmn.MediaType
	GetItem(string) (cmn.MediaItem, error)
}

type SearchResult struct {
	Next  string          `json:"next"`
	Items []cmn.MediaItem `json:"items"`
}
