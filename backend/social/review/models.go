package review

type ReviewRequest struct {
	ExternalID string `json:"external_id"`
	Rating     uint8  `json:"rating"`
	Body       string `json:"body"`
	Rewatch    bool   `json:"rewatch"`
	Date       string `json:"date"`
}

type ReviewResponse struct {
	ID         uint   `json:"id"`
	ExternalID string `json:"external_id"`
	Rating     uint8  `json:"rating"`
	Body       string `json:"body"`
	Rewatch    bool   `json:"rewatch"`
	Date       string `json:"date"`
}
