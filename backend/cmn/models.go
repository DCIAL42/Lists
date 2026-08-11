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
	UserID     string         `json:"user_id"`
	Status     TrackingStatus `json:"status"`
	Type       MediaType      `json:"type"`
	ExternalID string         `json:"external_id"`
}

type TrackingResponse struct {
	ID     uint           `json:"id,omitempty"`
	Status TrackingStatus `json:"status,omitempty"`
}

type MediaItem struct {
	Type       MediaType        `json:"type"`
	ExternalID string           `json:"external_id"`
	Data       any              `json:"data"`
	Tracking   TrackingResponse `json:"tracking,omitempty"`
}

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type Client interface {
	BuildURL(map[string]string) string
	TryRequest(context.Context, string) (*http.Response, error)
	ReadToSearchResult(*http.Response) (SearchResult, error)
	Search(context.Context, map[string]string) (SearchResult, error)
	GetItem(string) (MediaItem, error)
}

type SearchResult struct {
	Next  string      `json:"next"`
	Items []MediaItem `json:"items"`
}
