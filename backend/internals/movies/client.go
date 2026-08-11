package movies

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"sort"
	"strconv"

	"github.com/DCIAL42/media/cmn"
	"github.com/DCIAL42/media/internals/client"
)

func (c *Client) ReadToSearchResult(resp *http.Response) (res cmn.SearchResult, err error) {
	var data Response

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	movies := make([]Movie, 0, len(data.Results))

	for _, r := range data.Results {
		movies = append(movies, Movie{
			Title:      r.Title,
			Popularity: r.Popularity,
			Poster:     "https://image.tmdb.org/t/p/w500" + r.Poster,
		})
	}

	sort.Slice(movies, func(i, j int) bool {
		return movies[i].Popularity > movies[j].Popularity
	})

	results := make([]cmn.MediaItem, 0, len(movies))

	for i, m := range movies {
		results = append(results, cmn.MediaItem{
			Type:       cmn.TypeMovie,
			ExternalID: strconv.Itoa(data.Results[i].ExternalID),
			Data:       m,
		})
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

func NewMovieClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient,
		"https://api.themoviedb.org/3",
		"/search/movie",
		map[string]string{},
		map[string]string{
			"Authorization": "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIzYmFmYjliYjM1MjA2ZDQ0ZGNhYjU5MTc4NDNlYzdmNyIsIm5iZiI6MTc3MTI1OTA3OC40MjQsInN1YiI6IjY5OTM0NGM2MDZkZDViY2UxMzI4N2QzNSIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.Rx24Tn9as1XcsOypTv_Ufdx1_IWVAUfVL8RA4cIhBi0",
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
		Data: Movie{
			Title:      movie.Title,
			Popularity: movie.Popularity,
			Poster:     "https://image.tmdb.org/t/p/w500" + movie.Poster,
		},
	}, nil
}
