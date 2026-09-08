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

	"github.com/DCIAL42/lists/client"
	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
	"gorm.io/gorm"
)

func (r *AlbumSearchResponse) toAlbum() Album {
	var artist ArtistAPIResponse
	if len(r.Artists) > 0 {
		artist = r.Artists[0]
	}
	var cover string
	if len(r.Images) > 0 {
		cover = r.Images[0].URL
	}

	res := Album{
		Artist: artist.toArtist(),
		Media: cmn.Media{
			Type:       cmn.TypeAlbum,
			ExternalID: r.ExternalID,
			Name:       r.Name,
			Cover:      cover,
		},
	}

	return res
}

func (c *Client) ReadToSearchResult(resp *http.Response, userID string) (res cmn.SearchResult, err error) {
	var data SearchResponse

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		slog.Error(err.Error())
		return
	}
	externalIDs := make([]string, 0, len(data.Albums.Items))
	externalIDtoMovieResult := make(map[string]AlbumSearchResponse, len(data.Albums.Items))

	for _, r := range data.Albums.Items {
		externalIDs = append(externalIDs, r.ExternalID)
		externalIDtoMovieResult[r.ExternalID] = r
	}
	var albums []Album
	err = c.DB.Where("media_id IN (?)", c.DB.Model(&cmn.Media{}).Select("id").Where("external_id IN ?", externalIDs)).Preload("Media.Tracking").Preload("Media.Rating").Find(&albums).Error
	if err != nil {
		return
	}

	externalIDtoMovie := make(map[string]Album, len(albums))
	for _, a := range albums {
		externalIDtoMovie[a.Media.ExternalID] = a
	}

	results := make([]cmn.MediaResponse, 0, len(data.Albums.Items))

	for _, eID := range externalIDs {
		album, ok := externalIDtoMovie[eID]
		if !ok || album.ShouldUpdate() {
			newAlbum := externalIDtoMovieResult[eID]
			album = newAlbum.toAlbum()
			if _, err = db.TrySaveItem(c.DB, &album.Artist); err != nil {
				return
			}
			if _, err = db.TrySaveItem(c.DB, &album); err != nil {
				return
			}
		}
		results = append(results, album.toMediaResponse())
	}

	return cmn.SearchResult{Items: results}, nil
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
	params["fields"] = "albums(items(id,name,artists(id,name,images),images))"

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

	var tokenData struct {
		Token string `json:"access_token"`
		Type  string `json:"token_type"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenData)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	c.headers["Authorization"] = tokenData.Type + " " + tokenData.Token
}

func NewClient(httpClient *http.Client, DB *gorm.DB) *Client {
	c := &Client{
		DB,
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

func (c *Client) ResolveMedia(m cmn.Media) (res cmn.MediaResponse, err error) {
	switch m.Type {
	case cmn.TypeAlbum:
		c.FetchTracks(m.ID)
		var item Album
		result := c.DB.Where("media_id = ?", m.ID).Preload("Media").Preload("Tracks.Media").Preload("Artist.Media").First(&item)
		if result.Error != nil {
			err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "failed to get media"}
			return
		}
		item.Media = m
		err = c.FetchArtist(item.Artist.Media.ExternalID)
		return item.toMediaResponse(), nil
	case cmn.TypeArtist:
		var item Artist
		result := c.DB.Where("media_id = ?", m.ID).Preload("Media").Preload("Albums.Media").First(&item)
		if result.Error != nil {
			err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "failed to get media"}
			return
		}
		item.Media = m
		return item.toMediaResponse(), nil
	}
	return cmn.MediaResponse{}, &cmn.HttpError{Code: http.StatusInternalServerError, Message: "invalid type"}
}
