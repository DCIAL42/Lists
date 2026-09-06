package music

import (
	"net/http"

	"github.com/DCIAL42/lists/cmn"
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
	ExternalID string `json:"id"`
	Title      string `json:"name"`
	Artists    []struct {
		Name string `json:"name"`
	} `json:"artists"`
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
}

type SearchResponse struct {
	Albums struct {
		Items []AlbumSearchResponse `json:"items"`
		Next  string                `json:"next"`
	} `json:"albums"`
}

type TracksResponse struct {
	Items []struct {
		ExternalID string `json:"id"`
		Title      string `json:"name"`
		Duration   uint   `json:"duration_ms"`
	} `json:"items"`
}

type AlbumData struct {
	Artist string          `json:"artist"`
	Tracks []TrackResponse `json:"tracks"`
}

type Track struct {
	cmn.Model
	MediaID  uint
	Media    cmn.Media `gorm:"foreignKey:MediaID"`
	AlbumID  uint
	Duration uint
}

type TrackResponse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Duration uint   `json:"duration"`
}

type Album struct {
	cmn.Model
	MediaID uint
	Media   cmn.Media `gorm:"foreignKey:MediaID"`
	Artist  string
	Tracks  []Track `gorm:"foreignKey:AlbumID"`
}
