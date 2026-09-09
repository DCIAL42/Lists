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
	"github.com/DCIAL42/lists/db"
	"gorm.io/gorm"
)

func (r *MovieResponse) toMovie() Movie {
	return Movie{
		Popularity: r.Popularity,
		Media: cmn.Media{
			Type:       cmn.TypeMovie,
			ExternalID: strconv.Itoa(r.ExternalID),
			Name:       r.Title,
			Cover:      "https://image.tmdb.org/t/p/w500" + r.Poster,
		},
	}
}

func (r MovieResponse) GetExternalID() string {
	return strconv.Itoa(r.ExternalID)
}

func (r MovieResponse) ToExternalItem() *Movie {
	m := r.toMovie()
	return &m
}

func (r *Response) Items() []MovieResponse {
	return r.Results
}

func (m Movie) ToMediaResponse() (res cmn.MediaResponse) {
	res = cmn.MediaResponse{
		ID:    m.MediaID,
		Type:  cmn.TypeMovie,
		Name:  m.Media.Name,
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

func (m *Movie) GetMedia() *cmn.Media {
	return &m.Media
}

func (m *Movie) GetExternalID() string {
	return m.Media.ExternalID
}

func (m *Movie) GetID() uint {
	return m.ID
}

func (m *Movie) CacheItem(DB *gorm.DB) error {
	if _, err := db.TrySaveItem(DB, m); err != nil {
		return err
	}
	return nil
}

func (m *Movie) ShouldUpdate() bool {
	return db.DefaultShouldUpdate(m.Model.UpdatedAt)
}

func (c *Client) BuildURL(params map[string]string) string {
	query := url.Values{}

	for k, v := range params {
		query.Set(k, v)
	}

	url := c.baseURL + c.searchPath + "?" + query.Encode()
	return url
}

func (c *Client) DB() *gorm.DB {
	return c.db
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
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	for range 3 {
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

	return client.Search[*Movie, MovieResponse, *Response](ctx, c, params)
}

func (c *Client) ResolveMedia(m cmn.Media) (res cmn.MediaResponse, err error) {
	var item Movie
	result := c.db.Where("media_id = ?", m.ID).Preload("Media").First(&item)
	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "failed to get media"}
		return
	}
	item.Media = m
	return item.ToMediaResponse(), nil
}
