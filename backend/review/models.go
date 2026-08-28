package review

import (
	"time"

	"github.com/DCIAL42/lists/cmn"
)

type ReviewRequest struct {
	MediaID string `json:"media_id"`
	Rating  uint8  `json:"rating"`
	Body    string `json:"body"`
	Rewatch bool   `json:"rewatch"`
	Date    string `json:"date"`
}

type Review struct {
	cmn.Model
	UserID  string
	MediaID string
	Rating  uint8
	Body    string
	Rewatch bool
	Date    time.Time
}

type ReviewResponse struct {
	ID uint `json:"id"`
	ReviewRequest
}
