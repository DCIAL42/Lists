package music

import (
	"net/http"

	"github.com/DCIAL42/media/internals/client"
)

type Client struct {
	httpClient   *http.Client
	baseURL      string
	searchPath   string
	mediaType    client.MediaType
	configParams map[string]string
	headers      map[string]string
}

type Response struct {
	Albums struct {
		Items []AlbumResponse `json:"items"`
		Next  string          `json:"next"`
	} `json:"albums"`
}

type AlbumResponse struct {
	ExternalID string `json:"id"`
	Title      string `json:"name"`
	Artists    []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
}

type Album struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Cover  string `json:"cover"`
}
