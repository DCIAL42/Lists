package movies

import (
	"net/http"

	"github.com/DCIAL42/media/cmn"
)

type Client struct {
	httpClient   *http.Client
	baseURL      string
	searchPath   string
	mediaType    cmn.MediaType
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

type Movie struct {
	Title      string  `json:"title"`
	Popularity float32 `json:"popularity"`
	Poster     string  `json:"cover"`
}
