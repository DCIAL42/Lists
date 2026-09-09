package cmn

import (
	"context"
	"net/http"

	"gorm.io/gorm"
)

type ExternalItem interface {
	GetExternalID() string
	GetID() uint
	ShouldUpdate() bool
	GetMedia() *Media
	CacheItem(*gorm.DB) error
	ToMediaResponse() MediaResponse
}

type Client interface {
	BuildURL(map[string]string) string
	TryRequest(context.Context, string) (*http.Response, error)
	Search(ctx context.Context, params map[string]string) (SearchResult, error)
	ResolveMedia(Media) (MediaResponse, error)
	DB() *gorm.DB
}

type SearchResult struct {
	Next  string          `json:"next"`
	Items []MediaResponse `json:"items"`
}
