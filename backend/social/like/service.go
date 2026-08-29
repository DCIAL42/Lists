package like

import (
	"errors"
	"net/http"

	"github.com/DCIAL42/lists/cmn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	*gorm.DB
}

func NewService(DB *gorm.DB) *Service {
	return &Service{
		DB: DB,
	}
}

func (s *Service) GetLike(userID string, ListID uint) (res LikeResponse, err error) {
	var like Like
	err = s.DB.
		Where("user_id = ? AND list_id = ?", userID, ListID).
		First(&like).Error

	res.ID = like.ID

	if err == nil {
		res.Liked = true
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	return
}

func (s *Service) createLike(userID string, req LikeRequest) (res LikeResponse, err error) {
	var test Like

	result := s.DB.Unscoped().Where("list_id = ? AND user_id = ?", req.ListID, userID).First(&test)

	if result.Error == nil {
		result := s.DB.Unscoped().Model(&test).Update("deleted_at", nil)
		if result.Error != nil {
			err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
			return
		}
		if result.RowsAffected == 0 {
			err = &cmn.HttpError{Code: http.StatusNotFound, Message: "not found"}
			return
		}
		res = LikeResponse{
			ID:    test.ID,
			Liked: true,
		}
		return
	}

	if err := s.DB.
		Where("user_id = ? AND list_id = ?", userID, req.ListID).
		First(&Like{}).Error; err == nil {
		return LikeResponse{}, &cmn.HttpError{Code: http.StatusBadRequest, Message: "like already exists"}
	}

	like := Like{
		UserID: userID,
		ListID: req.ListID,
	}

	err = s.DB.Create(&like).Error

	res.ID = like.ID

	if err != nil {
		res.Liked = true
	}

	return
}

func (s *Service) deleteLike(userID string, listID uint) (res LikeResponse, err error) {
	var like Like
	result := s.DB.Clauses(clause.Returning{}).Where("list_id = ? AND user_id = ?", listID, userID).Delete(&like)

	if result.Error != nil {
		return res, result.Error
	}

	if result.RowsAffected == 0 {
		return res, &cmn.HttpError{Code: http.StatusNotFound, Message: "not found"}
	}
	res.ID = like.ID
	res.Liked = false

	return res, nil
}
