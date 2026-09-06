package cmn

import (
	"time"

	"gorm.io/gorm"
)

type MediaType string

const (
	TypeAlbum  MediaType = "album"
	TypeMovie  MediaType = "movie"
	TypeArtist MediaType = "artist"
	TypeTrack  MediaType = "track"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type Media struct {
	Model
	ExternalID string `gorm:"uniqueIndex;not null"`
	Type       MediaType
	Name       string
	Cover      string

	Tracking *TrackingItem `gorm:"foreignKey:MediaID"`
	Rating   *Rating       `gorm:"foreignKey:MediaID"`
}

type MediaResponse struct {
	ID       uint             `json:"id"`
	Type     MediaType        `json:"type"`
	Name     string           `json:"name"`
	Cover    string           `json:"cover"`
	Data     any              `json:"data"`
	Tracking TrackingResponse `json:"tracking"`
	Rating   RatingResponse   `json:"rating"`
}
