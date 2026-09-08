package movies

import (
	"context"
	"encoding/json"
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

func (m *Movie) ShouldUpdate() bool {
	return db.DefaultShouldUpdate(m.Model.UpdatedAt)
}

func (c *Client) ReadToSearchResult(resp *http.Response, userID string) (res cmn.SearchResult, err error) {
	var data Response

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	// sort.Slice(data.Results, func(i, j int) bool {
	// 	return data.Results[i].Popularity > data.Results[j].Popularity
	// })
	externalIDs := make([]string, 0, len(data.Results))
	externalIDtoMovieResult := make(map[string]MovieResponse, len(data.Results))

	for _, r := range data.Results {
		externalIDs = append(externalIDs, strconv.Itoa(r.ExternalID))
		externalIDtoMovieResult[strconv.Itoa(r.ExternalID)] = r
	}

	var movies []Movie
	err = c.DB.Where("media_id IN (?)", c.DB.Model(&cmn.Media{}).Select("id").Where("external_id IN ?", externalIDs)).Preload("Media.Tracking").Preload("Media.Rating").Find(&movies).Error
	if err != nil {
		return
	}

	externalIDtoMovie := make(map[string]Movie, len(movies))
	for _, movie := range movies {
		externalIDtoMovie[movie.Media.ExternalID] = movie
	}

	results := make([]cmn.MediaResponse, 0, len(data.Results))

	for _, eID := range externalIDs {
		movie, ok := externalIDtoMovie[eID]
		if !ok || movie.ShouldUpdate() {
			newMovie := externalIDtoMovieResult[eID]
			movie = newMovie.toMovie()
			if _, err = db.TrySaveItem(c.DB, &movie); err != nil {
				return
			}
		}
		results = append(results, movie.ToMediaResponse())
	}

	return cmn.SearchResult{Items: results}, nil
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

	return client.Search(ctx, c, params)
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
