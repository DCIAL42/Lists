package music

import "net/http"

type Client struct {
	httpClient   *http.Client
	baseURL      string
	searchPath   string
	resultType   string
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
	ExternalID string `json:"external_id"`
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Cover      string `json:"cover"`
	Type       string `json:"type"`
}
