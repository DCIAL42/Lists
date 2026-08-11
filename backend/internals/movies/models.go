package movies

import (
	"net/http"
)

type Client struct {
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

type Movie struct {
	Title      string  `json:"title"`
	Popularity float32 `json:"popularity"`
	Poster     string  `json:"cover"`
}
