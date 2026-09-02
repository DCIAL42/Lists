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
	"sort"
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
			Title:      r.Title,
			Cover:      "https://image.tmdb.org/t/p/w500" + r.Poster,
		},
	}
}

func (m *Movie) toMediaResponse() (res cmn.MediaResponse) {
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
		rating := *m.Media.Rating
		res.Rating = cmn.RatingResponse{
			ID:     rating.ID,
			Rating: rating.Rating,
		}
	}
	return
}

func (m Movie) GetID() uint {
	return m.ID
}

func (m Movie) GetExternalID() string {
	return m.Media.ExternalID
}

func (m Movie) GetModel() cmn.Model {
	return m.Model
}

func (c *Client) ReadToSearchResult(resp *http.Response, userID string) (res cmn.SearchResult, err error) {
	var data Response

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	sort.Slice(data.Results, func(i, j int) bool {
		return data.Results[i].Popularity > data.Results[j].Popularity
	})

	movies := make([]Movie, 0, len(data.Results))
	mediaIDs := make([]uint, 0, len(data.Results))

	for _, r := range data.Results {
		var movie Movie = r.toMovie()
		if _, err = db.TrySaveItem(c.DB, &movie); err != nil {
			return
		}
		movies = append(movies, movie)
		mediaIDs = append(mediaIDs, movie.MediaID)
	}

	var tracking []cmn.TrackingItem

	if err = c.DB.Where("user_id = ?", userID).Where("media_id IN ?", mediaIDs).Find(&tracking).Error; err != nil {
		return
	}

	trackingByMediaID := make(map[uint]*cmn.TrackingItem, len(tracking))
	for i := range tracking {
		trackingByMediaID[tracking[i].MediaID] = &tracking[i]
	}

	results := make([]cmn.MediaResponse, 0, len(movies))

	for i := range movies {
		movies[i].Media.Tracking = trackingByMediaID[movies[i].MediaID]
		results = append(results, movies[i].toMediaResponse())
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
	return item.toMediaResponse(), nil
}

func (c *Client) ResolveMedia(m cmn.Media) (res cmn.MediaResponse, err error) {
	var item Movie
	result := c.DB.Where("media_id = ?", m.ID).Preload("Media").First(&item)
	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "failed to get media"}
		return
	}
	item.Media = m
	return item.toMediaResponse(), nil
}
