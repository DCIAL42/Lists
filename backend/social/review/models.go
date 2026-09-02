package review

type ReviewRequest struct {
	MediaID uint    `json:"media_id"`
	Rating  float32 `json:"rating"`
	Body    string  `json:"body"`
	Rewatch bool    `json:"rewatch"`
	Date    string  `json:"date"`
}

type ReviewResponse struct {
	ID      uint    `json:"id"`
	MediaID uint    `json:"media_id"`
	Rating  float32 `json:"rating"`
	Body    string  `json:"body"`
	Rewatch bool    `json:"rewatch"`
	Date    string  `json:"date"`
}
