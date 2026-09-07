package cmn

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

	Media Media `gorm:"foreignKey:MediaID"`
}

func (t *TrackingItem) ToTrackingResponse() TrackingResponse {
	return TrackingResponse{
		ID:     t.ID,
		Status: t.Status,
	}
}

type TrackingResponse struct {
	ID     uint           `json:"id,omitempty"`
	Status TrackingStatus `json:"status,omitempty"`
}
