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

type TokenResponse struct {
	Token string `json:"access_token"`
	Type  string `json:"token_type"`
}

func (a *Album) toMediaResponse() (res cmn.MediaResponse) {
	res = cmn.MediaResponse{
		ID:    a.MediaID,
		Type:  a.Media.Type,
		Title: a.Media.Title,
		Cover: a.Media.Cover,
		Data: AlbumData{
			Artist: a.Artist,
		},
	}
	if a.Media.Tracking != nil {
		tracking := *a.Media.Tracking
		res.Tracking = cmn.TrackingResponse{
			ID:     tracking.ID,
			Status: tracking.Status,
		}
	}
	if a.Media.Rating != nil {
		rating := *a.Media.Rating
		res.Rating = cmn.RatingResponse{
			ID:     rating.ID,
			Rating: rating.Rating,
		}
	}
	return
}

func (r *AlbumResponse) toAlbum() Album {
	var artist string
	if len(r.Artists) > 0 {
		artist = r.Artists[0].Name
	}
	var cover string
	if len(r.Images) > 0 {
		cover = r.Images[0].URL
	}

	return Album{
		Artist: artist,
		Media: cmn.Media{
			Type:       cmn.TypeAlbum,
			ExternalID: r.ExternalID,
			Title:      r.Title,
			Cover:      cover,
		},
	}
}

func (a Album) GetID() uint {
	return a.ID
}

func (a Album) GetExternalID() string {
	return a.Media.ExternalID
}

func (a Album) GetModel() cmn.Model {
	return a.Model
}

func (c *Client) ReadToSearchResult(resp *http.Response, userID string) (res cmn.SearchResult, err error) {
	var data Response

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		slog.Error(err.Error())
		return
	}

	defer resp.Body.Close()

	albums := make([]Album, 0, len(data.Albums.Items))
	mediaIDs := make([]uint, 0, len(data.Albums.Items))

	for _, r := range data.Albums.Items {
		album := r.toAlbum()
		if _, err = db.TrySaveItem(c.DB, &album); err != nil {
			return
		}
		albums = append(albums, album)
		mediaIDs = append(mediaIDs, album.MediaID)
	}

	var tracking []cmn.TrackingItem

	if err = c.DB.Where("user_id = ?", userID).Where("media_id IN ?", mediaIDs).Find(&tracking).Error; err != nil {
		return
	}

	trackingByMediaID := make(map[uint]*cmn.TrackingItem, len(tracking))
	for i := range tracking {
		trackingByMediaID[tracking[i].MediaID] = &tracking[i]
	}

	results := make([]cmn.MediaResponse, 0, len(albums))

	for i := range albums {
		albums[i].Media.Tracking = trackingByMediaID[albums[i].MediaID]
		results = append(results, albums[i].toMediaResponse())
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
	params["fields"] = "albums(items(id,name,artists(name),images))"

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

func (c *Client) GetMedia(ID uint) (res cmn.MediaResponse, err error) {
	var item Album
	result := c.DB.Where("media_id = ?", ID).Preload("Media").First(&item)
	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "failed to get media"}
		return
	}
	return item.toMediaResponse(), nil
}

func (c *Client) ResolveMedia(m cmn.Media) (res cmn.MediaResponse, err error) {
	var item Album
	result := c.DB.Where("media_id = ?", m.ID).Preload("Media").First(&item)
	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "failed to get media"}
		return
	}
	item.Media = m
	return item.toMediaResponse(), nil
}
