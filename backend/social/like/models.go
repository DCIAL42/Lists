package like

import "github.com/DCIAL42/lists/cmn"

type LikeRequest struct {
	ListID uint `json:"list_id"`
}

type Like struct {
	cmn.Model
	UserID string
	ListID uint
}

type LikeResponse struct {
	ID    uint `json:"id,omitempty"`
	Liked bool `json:"liked"`
}
