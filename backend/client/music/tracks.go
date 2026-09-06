package music

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
)

type TracksResponse struct {
	Items []struct {
		ExternalID string `json:"id"`
		Name       string `json:"name"`
		Duration   uint   `json:"duration_ms"`
	} `json:"items"`
}

type Track struct {
	cmn.Model
	MediaID  uint
	Media    cmn.Media `gorm:"foreignKey:MediaID"`
	AlbumID  uint
	Duration uint
}

type TrackResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Duration uint   `json:"duration"`
}

func (t Track) GetExternalID() string {
	return t.Media.ExternalID
}
func (t Track) GetModel() cmn.Model {
	return t.Model
}
func (t Track) GetMedia() *cmn.Media {
	return &t.Media
}
func (t Track) GetMediaID() uint {
	return t.MediaID
}
func (t Track) ToMediaResponse() cmn.MediaResponse {
	return cmn.MediaResponse{}
}

func (c *Client) FetchTracks(albumMediaID uint) error {
	var track Track
	result := c.DB.Where("album_id = (?)", c.DB.Model(&Album{}).Select("id").Where("media_id = ?", albumMediaID)).First(&track)
	if result.Error == nil && time.Since(track.UpdatedAt) < 30*24*time.Hour {
		return nil
	}
	var item Album
	if err := c.DB.Where("media_id = ?", albumMediaID).Preload("Media").First(&item).Error; err != nil {
		return err
	}

	url := c.baseURL + "/albums/" + item.Media.ExternalID + "/tracks"
	resp, err := c.TryRequest(context.Background(), url)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	defer resp.Body.Close()

	var body TracksResponse

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		slog.Error(err.Error())
		return err
	}

	for _, t := range body.Items {
		track := Track{
			AlbumID:  item.ID,
			Duration: t.Duration,
			Media: cmn.Media{
				ExternalID: t.ExternalID,
				Name:       t.Name,
			},
		}
		_, err := db.TrySaveItem(c.DB, &track)
		if err != nil {
			return err
		}
	}
	return nil
}
