package tracking

import (
	"github.com/DCIAL42/media/cmn"
)

type TrackingItemResponse struct {
	ID uint `json:"id"`
	cmn.MediaItem
}

type TrackingListResponse struct {
	Items []cmn.MediaItem `json:"items"`
}
