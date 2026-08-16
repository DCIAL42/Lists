package music

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
	"strings"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/internals/client"
)

type TokenResponse struct {
	Token string `json:"access_token"`
	Type  string `json:"token_type"`
}

func (r *AlbumResponse) toMediaItem() cmn.MediaItem {
	var artist string
	if len(r.Artists) > 0 {
		artist = r.Artists[0].Name
	}

	var cover string
	if len(r.Images) > 0 {
		cover = r.Images[0].URL
	}

	return cmn.MediaItem{
		Type:       cmn.TypeAlbum,
		ExternalID: r.ExternalID,
		Data: Album{
			Title:  r.Title,
			Artist: artist,
			Cover:  cover,
		},
	}
}

func (c *Client) ReadToSearchResult(resp *http.Response) (res cmn.SearchResult, err error) {
	var data Response

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer resp.Body.Close()

	albums := make([]cmn.MediaItem, 0, len(data.Albums.Items))

	for _, r := range data.Albums.Items {
		albums = append(albums, r.toMediaItem())
	}

	return cmn.SearchResult{Items: albums}, nil
}

func (c *Client) BuildURL(params map[string]string) string {
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

	return c.baseURL + c.searchPath + "?" + query.Encode()
}

func (c *Client) fetchToken(ctx context.Context) {
	slog.Debug("Fetching api token.", "API", c.baseURL)
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

	resp, err := c.httpClient.Do(req)

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

	c.headers["Authorization"] = tokenData.Type + " " + tokenData.Token
}

func NewMusicClient(httpClient *http.Client) *Client {
	c := &Client{
		httpClient,
		"https://api.spotify.com/v1",
		"/search",
		map[string]string{
			"type":  "album",
			"limit": "10",
		},
		map[string]string{},
	}
	c.fetchToken(context.Background())
	return c
}

func (c *Client) TryRequest(ctx context.Context, targetUrl string) (*http.Response, error) {
	for range 3 {
		req, err := http.NewRequestWithContext(ctx, "GET", targetUrl, nil)

		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		for k, v := range c.headers {
			req.Header.Set(k, v)
		}

		resp, err := c.httpClient.Do(req)

		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}

		if resp.StatusCode == 401 {
			slog.Info("Response status not ok, refreshing token.", "StatusCode", resp.StatusCode, "API", c.baseURL)
			c.fetchToken(ctx)
			continue
		}

		return resp, nil
	}
	return nil, errors.New("Unable to make request")
}

func (c *Client) Search(ctx context.Context, params map[string]string) (cmn.SearchResult, error) {
	maps.Copy(params, c.configParams)

	return client.Search(ctx, c, params)
}

func (c *Client) GetItem(ID string) (res cmn.MediaItem, err error) {
	targetUrl := c.baseURL + "/albums/" + ID

	resp, err := c.TryRequest(context.Background(), targetUrl)

	if err != nil {
		return
	}

	var album AlbumResponse

	err = json.NewDecoder(resp.Body).Decode(&album)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer resp.Body.Close()

	return album.toMediaItem(), nil
}
