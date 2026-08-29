package like

import "github.com/DCIAL42/lists/cmn"

type LikeRequest struct {
	UserID string `json:"user_id"`
	ListID uint   `json:"list_id"`
}

type Like struct {
	cmn.Model
}
