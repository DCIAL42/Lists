package rating

import (
	"net/http"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
	"gorm.io/gorm"
)

type Service struct {
	*db.DBService
}

func NewService(DB *gorm.DB, clients map[cmn.MediaType]cmn.Client, config ...*db.ServiceConfig) *Service {
	return &Service{
		DBService: db.NewDBService(DB, clients, config...),
	}
}

func (s *Service) parseRatingRequest(req RatingRequest, userID string) (res cmn.Rating, err error) {
	if req.Rating > 10 {
		return cmn.Rating{}, &cmn.HttpError{Code: http.StatusBadRequest, Message: "invalid rating"}
	}
	if req.MediaID == 0 {
		return cmn.Rating{}, &cmn.HttpError{Code: http.StatusBadRequest, Message: "invalid media id"}
	}
	return cmn.Rating{
		UserID:  userID,
		MediaID: req.MediaID,
		Rating:  req.Rating,
	}, err
}

func toRatingResponse(r cmn.Rating) cmn.RatingResponse {
	return cmn.RatingResponse{
		ID:     r.ID,
		Rating: r.Rating,
	}
}

func (s *Service) createRating(userID string, req RatingRequest) (res cmn.RatingResponse, err error) {
	rating, err := s.parseRatingRequest(req, userID)

	if err != nil {
		err = &cmn.HttpError{Code: http.StatusNotImplemented, Message: "test"}
		return
	}

	result := s.DB.Create(&rating)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	res = toRatingResponse(rating)

	return
}

func (s *Service) updateRating(id uint, userID string, req RatingRequest) (res cmn.RatingResponse, err error) {
	result := s.DB.Model(&cmn.Rating{}).Where("id = ? AND user_id = ?", id, userID).Updates(req)

	if result.Error != nil {
		return cmn.RatingResponse{}, &cmn.HttpError{Code: http.StatusInternalServerError, Message: "failed to fetch rating"}
	}

	if result.RowsAffected == 0 {
		return cmn.RatingResponse{}, &cmn.HttpError{Code: http.StatusNotFound, Message: "not found"}
	}

	var rating cmn.Rating
	s.DB.First(&rating, id)

	return cmn.RatingResponse{
		ID:     rating.ID,
		Rating: rating.Rating,
	}, nil
}

func (s *Service) getAllRatings() (res []cmn.RatingResponse, err error) {
	ratings := make([]cmn.Rating, 0)
	result := s.DB.Find(&ratings)
	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	for _, rating := range ratings {
		res = append(res, toRatingResponse(rating))
	}

	return
}
