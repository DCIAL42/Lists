package movies

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"sort"

	"github.com/DCIAL42/media/internals/client"
)

func readToSearchResult(resp *http.Response) (client.SearchResult, error) {
	var data Response

	err := json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return client.SearchResult{}, err
	}

	movies := make([]Movie, 0, len(data.Results))

	for _, r := range data.Results {
		movies = append(movies, Movie{
			Id:         r.Id,
			Title:      r.Title,
			Popularity: r.Popularity,
			Poster:     "https://image.tmdb.org/t/p/w500" + r.Poster,
		})
	}

	sort.Slice(movies, func(i, j int) bool {
		return movies[i].Popularity > movies[j].Popularity
	})

	results := make([]any, 0, len(movies))

	for _, m := range movies {
		results = append(results, m)
	}

	return client.SearchResult{Items: results}, nil
}

func buildURL(baseURL string, params map[string]string) string {
	query := url.Values{}

	for k, v := range params {
		query.Set(k, v)
	}

	url := baseURL + "?" + query.Encode()
	return url
}

func fetchToken(ctx context.Context, c *client.Client) {}

func NewMovieClient(httpClient *http.Client) *client.Client {
	return client.NewClient(
		httpClient,
		"https://api.themoviedb.org/3/search/movie",
		"movie",
		map[string]string{},
		map[string]string{
			"Authorization": "Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiIzYmFmYjliYjM1MjA2ZDQ0ZGNhYjU5MTc4NDNlYzdmNyIsIm5iZiI6MTc3MTI1OTA3OC40MjQsInN1YiI6IjY5OTM0NGM2MDZkZDViY2UxMzI4N2QzNSIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.Rx24Tn9as1XcsOypTv_Ufdx1_IWVAUfVL8RA4cIhBi0",
			"Accept":        "application/json",
		},
		readToSearchResult,
		buildURL,
		fetchToken,
	)
}
