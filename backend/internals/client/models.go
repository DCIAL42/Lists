package client

import (
	"context"
	"net/http"
)

type MediaType string

const (
	TypeAlbum MediaType = "album"
	TypeMovie MediaType = "movie"
)

type Client interface {
	BuildURL(map[string]string) string
	TryRequest(context.Context, string) (*http.Response, error)
	ReadToSearchResult(*http.Response) (SearchResult, error)
	Search(context.Context, map[string]string) (SearchResult, error)
	GetMediaType() MediaType
	GetItem(string) (MediaItem, error)
}

type MediaItem struct {
	Type       MediaType `json:"type"`
	ExternalID string    `json:"external_id"`
	Data       any       `json:"data"`
}

type SearchResult struct {
	Next  string      `json:"next"`
	Items []MediaItem `json:"items"`
}
