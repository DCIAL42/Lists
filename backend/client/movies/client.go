package movies

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/DCIAL42/lists/client"
	"github.com/DCIAL42/lists/cmn"
	"gorm.io/gorm"
)

func (r *MovieResponse) toMovie() *Movie {
	return &Movie{
		Popularity: r.Popularity,
		Media: cmn.Media{
			Type:       cmn.TypeMovie,
			ExternalID: strconv.Itoa(r.ExternalID),
			Title:      r.Title,
			Cover:      "https://image.tmdb.org/t/p/w500" + r.Poster,
		},
	}
}

func (m Movie) ToMediaResponse() (res cmn.MediaResponse) {
	res = cmn.MediaResponse{
		ID:    m.MediaID,
		Type:  cmn.TypeMovie,
		Title: m.Media.Title,
		Cover: m.Media.Cover,
		Data: MovieData{
			Popularity: m.Popularity,
		},
	}
	if m.Media.Tracking != nil {
		tracking := *m.Media.Tracking
		res.Tracking = cmn.TrackingResponse{
			ID:     tracking.ID,
			Status: tracking.Status,
		}
	}
	if m.Media.Rating != nil {
		res.Rating = (*m.Media.Rating).ToRatingResponse()
	}
	return
}

func (m Movie) GetID() uint {
	return m.ID
}

func (m Movie) GetMediaID() uint {
	return m.MediaID
}

func (m Movie) GetMedia() *cmn.Media {
	return &m.Media
}

func (m Movie) GetExternalID() string {
	return m.Media.ExternalID
}

func (m Movie) GetModel() cmn.Model {
	return m.Model
}

func (r Response) Items() []MovieResponse {
	return r.Results
}

func (m MovieResponse) ToDBItem() *Movie {
	return m.toMovie()
}

func (c *Client) ReadToSearchResult(resp *http.Response, userID string) (res cmn.SearchResult, err error) {
	return client.TestRead[*Movie, MovieResponse, Response](c.DB, resp, userID)
}

func (c *Client) BuildURL(params map[string]string) string {
	query := url.Values{}

	for k, v := range params {
		query.Set(k, v)
	}

	url := c.baseURL + c.searchPath + "?" + query.Encode()
	return url
}

func NewClient(httpClient *http.Client, DB *gorm.DB) *Client {
	token, ok := os.LookupEnv("TMDB_TOKEN")
	if !ok {
		panic("tmdb token not found")
	}
	return &Client{
		DB,
		httpClient,
		"https://api.themoviedb.org/3",
		"/search/movie",
		map[string]string{},
		map[string]string{
			"Authorization": "Bearer " + token,
			"Accept":        "application/json",
		},
	}
}

func (c *Client) TryRequest(ctx context.Context, url string) (*http.Response, error) {
	for range 3 {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		for k, v := range c.headers {
			req.Header.Set(k, v)
		}

		res, err := c.httpClient.Do(req)

		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		if res.StatusCode != 200 {
			slog.Error("Response not ok, trying again.", "StatusCode", res.StatusCode, "API", c.baseURL)
			continue
		}

		return res, nil
	}
	return nil, errors.New("Unable to make request")
}

func (c *Client) Search(ctx context.Context, params map[string]string) (cmn.SearchResult, error) {
	maps.Copy(params, c.configParams)

	return client.Search(ctx, c, params)
}

func (c *Client) GetMedia(ID uint) (res cmn.MediaResponse, err error) {
	var item Movie
	result := c.DB.Where("media_id = ?", ID).Preload("Media").First(&item)
	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "failed to get media"}
		return
	}
	return item.ToMediaResponse(), nil
}

func (c *Client) ResolveMedia(m cmn.Media) (res cmn.MediaResponse, err error) {
	var item Movie
	result := c.DB.Where("media_id = ?", m.ID).Preload("Media").First(&item)
	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "failed to get media"}
		return
	}
	item.Media = m
	return item.ToMediaResponse(), nil
}
