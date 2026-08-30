package like

type LikeRequest struct {
	ListID uint `json:"list_id"`
}

type LikeResponse struct {
	ID    uint `json:"id,omitempty"`
	Liked bool `json:"liked"`
}
