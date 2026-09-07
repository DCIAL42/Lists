package music

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
)

type ArtistAPIResponse struct {
	ExternalID string `json:"id"`
	Name       string `json:"name"`
	Images     []struct {
		URL string `json:"url"`
	} `json:"images"`
}

type Artist struct {
	cmn.Model
	MediaID uint
	Media   cmn.Media `gorm:"foreignKey:MediaID"`

	Albums []Album `gorm:"foreignKey:ArtistID"`
}

type ArtistData struct {
	Albums []cmn.MediaResponse `json:"albums"`
}

type ArtistResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Cover string `json:"cover"`
}

func (a *ArtistAPIResponse) toArtist() Artist {
	cover := ""
	if len(a.Images) > 0 {
		cover = a.Images[0].URL
	}

	return Artist{
		Media: cmn.Media{
			ExternalID: a.ExternalID,
			Type:       cmn.TypeArtist,
			Name:       a.Name,
			Cover:      cover,
		},
	}
}

func (a Artist) GetExternalID() string {
	return a.Media.ExternalID
}
func (a Artist) GetID() uint {
	return a.ID
}
func (a Artist) ShouldUpdate() bool {
	return db.DefaultShouldUpdate(a.Media.UpdatedAt)
}
func (a Artist) GetMedia() *cmn.Media {
	return &a.Media
}
func (a Artist) toMediaResponse() cmn.MediaResponse {
	albums := make([]cmn.MediaResponse, 0)
	for _, album := range a.Albums {
		albums = append(albums, album.toMediaResponse())
	}
	res := cmn.MediaResponse{
		ID:    a.Media.ID,
		Type:  a.Media.Type,
		Name:  a.Media.Name,
		Cover: a.Media.Cover,
		Data: ArtistData{
			Albums: albums,
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
		res.Rating = (*a.Media.Rating).ToRatingResponse()
	}
	return res
}

func (c *Client) FetchArtist(externalID string) error {
	if externalID == "" {
		return &cmn.HttpError{Code: http.StatusInternalServerError, Message: "invalid artist external id"}
	}
	var artist Artist
	result := c.DB.Where("media_id = (?)", c.DB.Model(&cmn.Media{}).Select("id").Where("external_id = ?", externalID)).First(&artist)
	if result.Error == nil && time.Since(artist.UpdatedAt) < time.Second {
		return nil
	}

	url := c.baseURL + "/artists/" + externalID
	resp, err := c.TryRequest(context.Background(), url)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	defer resp.Body.Close()

	var body ArtistAPIResponse

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		slog.Error(err.Error())
		return err
	}

	cover := ""
	if len(body.Images) > 0 {
		cover = body.Images[0].URL
	}

	if result.Error == nil {
		artist = Artist{
			Model: cmn.Model{
				ID: artist.ID,
			},
			Media: cmn.Media{
				Model: cmn.Model{
					ID: artist.MediaID,
				},
				ExternalID: externalID,
				Type:       cmn.TypeArtist,
				Name:       body.Name,
				Cover:      cover,
			},
		}
	} else {
		artist = Artist{
			Media: cmn.Media{
				ExternalID: externalID,
				Type:       cmn.TypeArtist,
				Name:       body.Name,
				Cover:      cover,
			},
		}
	}

	_, err = db.TrySaveItem(c.DB, &artist)
	return err
}
