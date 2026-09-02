package movies

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

type Response struct {
	Results []MovieResponse `json:"results"`
}

type MovieResponse struct {
	ExternalID int     `json:"id"`
	Title      string  `json:"title"`
	Popularity float32 `json:"popularity"`
	Poster     string  `json:"poster_path"`
}

type MovieData struct {
	Popularity float32 `json:"popularity"`
}

type Movie struct {
	cmn.Model
	MediaID    uint      `gorm:"uniqueIndex;not null"`
	Media      cmn.Media `gorm:"foreignKey:MediaID"`
	Popularity float32
}
