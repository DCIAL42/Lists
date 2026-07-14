package music

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/DCIAL42/media/internals/client"
)

type TokenResponse struct {
	Token string `json:"access_token"`
	Type  string `json:"token_type"`
}

func readToSearchResult(resp *http.Response) (client.SearchResult, error) {
	var data Response

	err := json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return client.SearchResult{}, err
	}

	albums := make([]any, 0, len(data.Albums.Items))

	for _, r := range data.Albums.Items {
		var artist string
		if len(r.Artists) > 0 {
			artist = r.Artists[0].Name
		}

		var cover string
		if len(r.Images) > 0 {
			cover = r.Images[0].URL
		}

		albums = append(albums, Album{
			Id:     r.Id,
			Title:  r.Title,
			Artist: artist,
			Cover:  cover,
		})
	}

	return client.SearchResult{Items: albums}, nil
}

func buildURL(baseURL string, params map[string]string) string {
	query := url.Values{}

	params["q"] = params["query"]
	delete(params, "query")

	page, err := strconv.Atoi(params["page"])
	if err != nil {
		page = 0
	}
	params["offset"] = strconv.Itoa(page * 10)

	for k, v := range params {
		query.Set(k, v)
	}

	url := baseURL + "?" + query.Encode()
	return url
}

func fetchToken(ctx context.Context, c *client.Client) {
	id := os.Getenv("SPOTIFY_CLIENT_ID")
	secret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", id)
	values.Set("client_secret", secret)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", strings.NewReader(values.Encode()))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.GetClient().Do(req)

	if err != nil || resp.StatusCode != 200 {
		slog.Error(err.Error())
		return
	}
	defer resp.Body.Close()

	var tokenData TokenResponse

	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	c.SetHeader("Authorization", fmt.Sprintf("%s %s", tokenData.Type, tokenData.Token))
}

func NewMusicClient(httpClient *http.Client) *client.Client {
	return client.NewClient(
		httpClient,
		"https://api.spotify.com/v1/search",
		"album",
		map[string]string{
			"type":  "album",
			"limit": "10",
		},
		map[string]string{},
		readToSearchResult,
		buildURL,
		fetchToken,
	)
}
