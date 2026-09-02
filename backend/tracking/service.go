package tracking

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

func (s *Service) getTrackingItem(req cmn.TrackingItem) (res cmn.TrackingResponse, err error) {
	var item cmn.TrackingItem
	result := s.DB.Where(req).First(&item)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	return cmn.TrackingResponse{
		ID:     item.ID,
		Status: item.Status,
	}, nil
}

func (s *Service) createTrackingItem(req cmn.TrackingItem) (res cmn.TrackingResponse, err error) {
	result := s.DB.Create(&req)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	res = cmn.TrackingResponse{
		ID:     req.ID,
		Status: req.Status,
	}

	return
}

func (s *Service) updateTrackingItem(req cmn.TrackingItem) (res cmn.TrackingResponse, err error) {
	result := s.DB.Where("id = ? AND user_id = ?", req.ID, req.UserID).Updates(req)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	if result.RowsAffected == 0 {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: "not found"}
		return
	}

	var item cmn.TrackingItem
	s.DB.First(&item, req.ID)

	return cmn.TrackingResponse{
		ID:     item.ID,
		Status: item.Status,
	}, nil
}

func (s *Service) deleteTrackingItem(id uint, userID string) (res cmn.TrackingResponse, err error) {
	var item cmn.TrackingItem
	result := s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&item)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	if result.RowsAffected == 0 {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: "not found"}
		return
	}

	return cmn.TrackingResponse{
		ID:     item.ID,
		Status: item.Status,
	}, nil
}

func (s *Service) getTrackingList(pat TrackingItemQuery, page int) (res TrackingListResponse, err error) {
	list := make([]cmn.TrackingItem, 0)

	var count int64
	result := s.DB.
		Model(&cmn.TrackingItem{}).
		Preload("Media").
		Where(
			"user_id = ? AND status = ? AND type IN ?",
			pat.UserID,
			pat.Status,
			pat.Types,
		).
		Order("id desc").
		Count(&count).
		Offset((int(page) - 1) * int(s.PageSize)).
		Limit(int(s.PageSize)).
		Find(&list)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: "No items in tracking list"}
		return
	}

	resolved := make([]cmn.MediaResponse, 0, len(list))

	for _, item := range list {
		var resItem cmn.MediaResponse
		resItem, err = s.ResolveMedia(item.MediaID, pat.UserID)

		if err != nil {
			err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "Error with api"}
			return
		}

		resItem.Tracking = cmn.TrackingResponse{
			ID:     item.ID,
			Status: item.Status,
		}

		resolved = append(resolved, resItem)
	}

	res = TrackingListResponse{
		Items: resolved,
		Count: int(count),
	}

	return
}

func (s *Service) getAllTrackingItems() (res []cmn.TrackingItem, err error) {
	res = make([]cmn.TrackingItem, 0)

	result := s.DB.Find(&res)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: "No items in tracking list"}
		return
	}

	return
}
