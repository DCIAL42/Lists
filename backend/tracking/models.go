package tracking

import (
	"time"

	"github.com/DCIAL42/media/cmn"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

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
	Type       cmn.MediaType  `json:"type"`
	ExternalID string         `json:"external_id"`
}

type TrackingItemResponse struct {
	ID uint `json:"id"`
	cmn.MediaItem
}

type TrackingListResponse struct {
	Items []TrackingItemResponse `json:"items"`
}
