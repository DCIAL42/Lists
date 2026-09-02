package cmn

import "time"

type Rating struct {
	Model
	UserID  string  `json:"user_id"`
	MediaID uint    `json:"media_id"`
	Rating  float32 `json:"rating"`
}

type RatingResponse struct {
	ID     uint    `json:"id,omitempty"`
	Rating float32 `json:"rating,omitempty"`
}

type Review struct {
	Model
	UserID  string
	MediaID uint
	Rating  float32
	Body    string
	Rewatch bool
	Date    time.Time

	Media Media `gorm:"foreignKey:MediaID"`
}

type Like struct {
	Model
	UserID string
	ListID uint
}
