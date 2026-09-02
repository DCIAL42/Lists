package rating

type RatingRequest struct {
	MediaID uint    `json:"media_id"`
	Rating  float32 `json:"rating"`
}
