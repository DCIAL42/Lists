package cmn

import (
	"context"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type MediaType string

const (
	TypeAlbum MediaType = "album"
	TypeMovie MediaType = "movie"
)

type TrackingStatus string

const (
	Backlog   TrackingStatus = "backlog"
	Paused    TrackingStatus = "paused"
	Completed TrackingStatus = "completed"
)

type TrackingItem struct {
	Model
	UserID  string         `json:"user_id"`
	Status  TrackingStatus `json:"status"`
	Type    MediaType      `json:"type"`
	MediaID uint           `json:"media_id"`

	Media Media `gorm:"foreignKey:MediaID"`
}

type TrackingResponse struct {
	ID     uint           `json:"id,omitempty"`
	Status TrackingStatus `json:"status,omitempty"`
}

type Rating struct {
	Model
	UserID  string `json:"user_id"`
	MediaID uint   `json:"media_id"`
	Rating  uint8  `json:"rating"`
}

type RatingResponse struct {
	ID     uint  `json:"id,omitempty"`
	Rating uint8 `json:"rating"`
}

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type MediaItem struct {
	Type       MediaType        `json:"type"`
	ExternalID string           `json:"external_id"`
	Data       any              `json:"data"`
	Tracking   TrackingResponse `json:"tracking,omitempty"`
	Rating     RatingResponse   `json:"rating,omitempty"`
	Cover      string           `json:"cover"`
}

type Review struct {
	Model
	UserID  string
	MediaID uint
	Rating  uint8
	Body    string
	Rewatch bool
	Date    time.Time

	Media Media `gorm:"foreignKey:MediaID"`
}

type Media struct {
	Model
	ExternalID string `gorm:"uniqueIndex;not null"`
	Type       MediaType
	Title      string
	Cover      string
	Reviews    []Review `gorm:"foreignKey:MediaID"`
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

type Client interface {
	BuildURL(map[string]string) string
	TryRequest(context.Context, string) (*http.Response, error)
	ReadToSearchResult(*http.Response) (SearchResult, error)
	Search(context.Context, map[string]string) (SearchResult, error)
	GetItem(string) (MediaItem, error)
	GetMedia(uint) (MediaResponse, error)
}

type SearchResult struct {
	Next  string          `json:"next"`
	Items []MediaResponse `json:"items"`
}

type Like struct {
	Model
	UserID string
	ListID uint
}
