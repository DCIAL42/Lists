package cmn

import (
	"time"

	"gorm.io/gorm"
)

type MediaType string

const (
	TypeAlbum MediaType = "album"
	TypeMovie MediaType = "movie"
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
	Title      string
	Cover      string

	Tracking *TrackingItem `gorm:"foreignKey:MediaID"`
	Rating   *Rating       `gorm:"foreignKey:MediaID"`
}

type MediaResponse struct {
	ID       uint             `json:"id"`
	Type     MediaType        `json:"type"`
	Title    string           `json:"title"`
	Cover    string           `json:"cover"`
	Data     any              `json:"data"`
	Tracking TrackingResponse `json:"tracking"`
	Rating   RatingResponse   `json:"rating"`
}
