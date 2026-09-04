package cmn

import "time"

type Rating struct {
	Model
	UserID  string `gorm:"not null;uniqueIndex:idx_user_rating"`
	MediaID uint   `gorm:"not null;uniqueIndex:idx_user_rating"`
	Rating  float32
}

type RatingResponse struct {
	ID     uint    `json:"id,omitempty"`
	Rating float32 `json:"rating,omitempty"`
}

func (r *Rating) ToRatingResponse() RatingResponse {
	return RatingResponse{
		ID:     r.ID,
		Rating: r.Rating,
	}
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
