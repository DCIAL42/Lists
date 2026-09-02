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
	UserID  string         `json:"user_id" gorm:"not null;uniqueIndex:idx_user_tracking"`
	MediaID uint           `json:"media_id" gorm:"not null;uniqueIndex:idx_user_tracking"`
	Status  TrackingStatus `json:"status"`
	Type    MediaType      `json:"type"`

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

	Tracking *TrackingItem `gorm:"foreignKey:MediaID"`

	Reviews []Review `gorm:"foreignKey:MediaID"`
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
	ReadToSearchResult(*http.Response, string) (SearchResult, error)
	Search(ctx context.Context, params map[string]string) (SearchResult, error)
	GetMedia(uint) (MediaResponse, error)
	ResolveMedia(Media) (MediaResponse, error)
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
