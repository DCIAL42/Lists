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

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
	"github.com/DCIAL42/lists/internals/client"
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

func (r *MovieResponse) toMediaItem() cmn.MediaItem {
	return cmn.MediaItem{
		Type:       cmn.TypeMovie,
		ExternalID: strconv.Itoa(r.ExternalID),
		Cover:      "https://image.tmdb.org/t/p/w500" + r.Poster,
		Data: MovieData{
			Title:      r.Title,
			Popularity: r.Popularity,
		},
	}
}

func (m *Movie) toMediaResponse() cmn.MediaResponse {
	return cmn.MediaResponse{
		ID:    m.MediaID,
		Type:  cmn.TypeMovie,
		Title: m.Media.Title,
		Cover: m.Media.Cover,
		Data: MovieData{
			Popularity: m.Popularity,
		},
	}
}

func (m *Movie) toMediaItem() cmn.MediaItem {
	return cmn.MediaItem{
		Type:       cmn.TypeMovie,
		ExternalID: m.Media.ExternalID,
		Cover:      m.Media.Cover,
		Data: MovieData{
			Title:      m.Media.Title,
			Popularity: m.Popularity,
		},
	}
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

func (c *Client) ReadToSearchResult(resp *http.Response) (res cmn.SearchResult, err error) {
	var data Response

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	sort.Slice(data.Results, func(i, j int) bool {
		return data.Results[i].Popularity > data.Results[j].Popularity
	})

	results := make([]cmn.MediaResponse, 0, len(data.Results))

	for _, r := range data.Results {
		var movie Movie = r.toMovie()
		db.TrySaveItem(c.DB, &movie)
		results = append(results, movie.toMediaResponse())
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

func NewMovieClient(httpClient *http.Client, DB *gorm.DB) *Client {
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

func (c *Client) GetItem(ID string) (res cmn.MediaItem, err error) {
	var item Movie
	ok := db.TryGetItem(c.DB, ID, &item)
	if ok {
		return item.toMediaItem(), nil
	}
	if len(ID) == 0 {
		err = errors.ErrUnsupported
		return
	}
	targetUrl := c.baseURL + "/movie/" + ID

	resp, err := c.TryRequest(context.Background(), targetUrl)

	if err != nil {
		return
	}

	var movie MovieResponse

	err = json.NewDecoder(resp.Body).Decode(&movie)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer resp.Body.Close()

	return cmn.MediaItem{
		Type:       cmn.TypeMovie,
		ExternalID: strconv.Itoa(movie.ExternalID),
		Cover:      "https://image.tmdb.org/t/p/w500" + movie.Poster,
		Data: MovieData{
			Title:      movie.Title,
			Popularity: movie.Popularity,
		},
	}, nil
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
