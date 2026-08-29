package music

import (
	"net/http"

	"github.com/DCIAL42/lists/db"
)

type Service struct {
	*db.DBService
}

type Client struct {
	httpClient   *http.Client
	baseURL      string
	searchPath   string
	configParams map[string]string
	headers      map[string]string
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

type Response struct {
	Albums struct {
		Items []AlbumResponse `json:"items"`
		Next  string          `json:"next"`
	} `json:"albums"`
}

type AlbumData struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Cover  string `json:"cover"`
}

type Album struct {
	ID string `gorm:"primaryKey"`
	AlbumData
}
