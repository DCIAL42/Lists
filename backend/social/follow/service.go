package follow

import (
	"errors"
	"net/http"

	"github.com/DCIAL42/lists/cmn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) newFollow(follower, followed string) (res FollowResponse, err error) {
	var test Follow

	result := s.DB.Unscoped().Where("follower = ? AND followed = ?", follower, followed).First(&test)

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
		res = FollowResponse{
			ID:       test.ID,
			Followed: true,
		}
		return
	}
	f := Follow{Follower: follower, Followed: followed}
	err = s.DB.Create(&f).Error

	return FollowResponse{ID: f.ID, Followed: err == nil}, err
}

func (s *Service) getFollow(follower, followed string) (res FollowResponse, err error) {
	var f Follow
	if err := s.DB.Where("follower = ? AND followed = ?", follower, followed).First(&f).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return FollowResponse{}, nil
		}
		return FollowResponse{}, &cmn.HttpError{Code: http.StatusInternalServerError, Message: "not finished"}
	}
	return FollowResponse{f.ID, true}, nil
}

func (s *Service) deleteFollow(userID string, followID uint) (res FollowResponse, err error) {
	var f Follow

	if err := s.DB.Clauses(clause.Returning{}).
		Where("follower = ? AND id = ?", userID, followID).
		Delete(&f).Error; err != nil {
		return FollowResponse{}, &cmn.HttpError{Code: http.StatusInternalServerError, Message: "not finished"}
	}

	return FollowResponse{f.ID, false}, nil
}

func (s *Service) getFollowers(followedID string) (res []FollowResponse, err error) {
	follows := make([]Follow, 0)
	result := s.DB.Where("followed = ?", followedID).Find(&follows)
	res = make([]FollowResponse, len(follows))
	for _, r := range follows {
		res = append(res, FollowResponse{
			r.ID,
			true,
		})
	}
	return res, result.Error
}
