package tracking

import (
	"github.com/DCIAL42/lists/cmn"
)

type TrackingItemQuery struct {
	UserID string
	Status cmn.TrackingStatus
	Types  []cmn.MediaType
}

type TrackingListResponse struct {
	Count int                 `json:"count"`
	Items []cmn.MediaResponse `json:"items"`
}
