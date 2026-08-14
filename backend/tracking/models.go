package tracking

import (
	"github.com/DCIAL42/media/cmn"
)

type TrackingItemQuery struct {
	UserID string
	Status cmn.TrackingStatus
	Types  []cmn.MediaType
}

type TrackingItemResponse struct {
	ID uint `json:"id"`
	cmn.MediaItem
}

type TrackingListResponse struct {
	Items []cmn.MediaItem `json:"items"`
}
