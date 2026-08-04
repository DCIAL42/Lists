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

type MediaItem struct {
	Type       MediaType `json:"type"`
	ExternalID string    `json:"external_id"`
	Data       any       `json:"data"`
}

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}
