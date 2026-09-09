package music

import (
	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
	"gorm.io/gorm"
)

type AlbumData struct {
	Artist ArtistResponse  `json:"artist"`
	Tracks []TrackResponse `json:"tracks"`
}

type Album struct {
	cmn.Model
	MediaID  uint
	Media    cmn.Media `gorm:"foreignKey:MediaID"`
	ArtistID uint
	Artist   Artist  `gorm:"foreignKey:ArtistID"`
	Tracks   []Track `gorm:"foreignKey:AlbumID"`
}

func (a Album) ToMediaResponse() (res cmn.MediaResponse) {
	tracks := make([]TrackResponse, 0, len(a.Tracks))
	for _, track := range a.Tracks {
		tracks = append(tracks, TrackResponse{
			ID:       track.ID,
			Name:     track.Media.Name,
			Duration: track.Duration,
		})
	}
	res = cmn.MediaResponse{
		ID:    a.MediaID,
		Type:  a.Media.Type,
		Name:  a.Media.Name,
		Cover: a.Media.Cover,
		Data: AlbumData{
			Artist: ArtistResponse{
				ID:    a.Artist.MediaID,
				Name:  a.Artist.Media.Name,
				Cover: a.Artist.Media.Cover,
			},
			Tracks: tracks,
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
	return
}

func (a Album) GetExternalID() string {
	return a.Media.ExternalID
}

func (a Album) GetMedia() *cmn.Media {
	return &a.Media
}

func (a *Album) GetID() uint {
	return a.ID
}

func (a *Album) ShouldUpdate() bool {
	return db.DefaultShouldUpdate(a.Model.UpdatedAt)
}

func (a *Album) CacheItem(DB *gorm.DB) error {
	if _, err := db.TrySaveItem(DB, &a.Artist); err != nil {
		return err
	}
	if _, err := db.TrySaveItem(DB, a); err != nil {
		return err
	}
	return nil
}
