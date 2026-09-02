package review

import (
	"net/http"
	"time"

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

func (s *Service) parseReviewRequest(req ReviewRequest, userID string) (cmn.Review, error) {
	if req.Rating > 10 {
		return cmn.Review{}, &cmn.HttpError{Code: http.StatusBadRequest, Message: "invalid rating"}
	}
	if req.MediaID == 0 {
		return cmn.Review{}, &cmn.HttpError{Code: http.StatusBadRequest, Message: "invalid media id"}
	}
	d, err := time.Parse("02-01-2006", req.Date)
	if err != nil {
		return cmn.Review{}, &cmn.HttpError{Code: http.StatusBadRequest, Message: "invalid date"}
	}
	media, err := s.GetMedia(req.MediaID, userID)
	return cmn.Review{
		UserID:  userID,
		MediaID: media.ID,
		Rating:  req.Rating,
		Body:    req.Body,
		Rewatch: req.Rewatch,
		Date:    d,
	}, err
}

func toReviewResponse(r cmn.Review) ReviewResponse {
	return ReviewResponse{
		ID:      r.ID,
		MediaID: r.MediaID,
		Rating:  r.Rating,
		Body:    r.Body,
		Rewatch: r.Rewatch,
		Date:    r.Date.Format("02-01-2006"),
	}
}

func (s *Service) createReview(userID string, req ReviewRequest) (res ReviewResponse, err error) {
	review, err := s.parseReviewRequest(req, userID)

	if err != nil {
		err = &cmn.HttpError{Code: http.StatusNotImplemented, Message: "test"}
		return
	}

	result := s.DB.Create(&review)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	res = toReviewResponse(review)

	return
}

func (s *Service) getAllReviews() (res []ReviewResponse, err error) {
	reviews := make([]cmn.Review, 0)
	result := s.DB.Find(&reviews)
	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	for _, review := range reviews {
		res = append(res, toReviewResponse(review))
	}

	return
}

func (s *Service) getReviewsByUser(userID string) ([]ReviewResponse, error) {
	var reviews []cmn.Review
	if err := s.DB.Where("user_id = ?", userID).Find(&reviews).Error; err != nil {
		return nil, err
	}

	res := make([]ReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		res = append(res, toReviewResponse(r))
	}

	return res, nil
}

func (s *Service) getReviewsByMedia(mediaID uint) ([]ReviewResponse, error) {
	var reviews []cmn.Review
	if err := s.DB.Where("media_id = ?", mediaID).Find(&reviews).Error; err != nil {
		return nil, &cmn.HttpError{Code: http.StatusInternalServerError, Message: "failed to find reviews"}
	}

	res := make([]ReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		res = append(res, toReviewResponse(r))
	}

	return res, nil
}
