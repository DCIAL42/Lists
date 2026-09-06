package music

import (
	"net/http"

	"gorm.io/gorm"
)

type Client struct {
	*gorm.DB
	httpClient   *http.Client
	baseURL      string
	searchPath   string
	configParams map[string]string
	headers      map[string]string
}

type AlbumSearchResponse struct {
	ExternalID string              `json:"id"`
	Name       string              `json:"name"`
	Artists    []ArtistAPIResponse `json:"artists"`
	Images     []struct {
		URL string `json:"url"`
	} `json:"images"`
}

type SearchResponse struct {
	Albums struct {
		Items []AlbumSearchResponse `json:"items"`
		Next  string                `json:"next"`
	} `json:"albums"`
}
