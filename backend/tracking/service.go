package tracking

import (
	"net/http"

	"github.com/DCIAL42/media/cmn"
	"github.com/DCIAL42/media/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	*db.DBService
}

func NewService(DB *gorm.DB, clients map[cmn.MediaType]cmn.Client, config ...*db.ServiceConfig) *Service {
	return &Service{
		DBService: db.NewDBService(DB, clients, config...),
	}
}

func (s *Service) GetTrackingItem(req cmn.TrackingItem) (res cmn.TrackingItem, err error) {
	result := s.DB.Where(req).First(&res)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	return
}

func (s *Service) createTrackingItem(req cmn.TrackingItem) (res cmn.TrackingItem, err error) {
	result := s.DB.Create(&req)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	res = req

	return
}

func (s *Service) updateTrackingItem(req cmn.TrackingItem) (res cmn.TrackingItem, err error) {
	result := s.DB.Where("id = ? AND user_id = ?", req.ID, req.UserID).Updates(req)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	if result.RowsAffected == 0 {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: "not found"}
		return
	}

	s.DB.First(&res, req.ID)

	return
}

func (s *Service) deleteTrackingItem(id uint, userID string) (res cmn.TrackingItem, err error) {
	result := s.DB.Clauses(clause.Returning{}).Where("id = ? AND user_id = ?", id, userID).Delete(&res)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	if result.RowsAffected == 0 {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: "not found"}
		return
	}

	return
}

func (s *Service) getTrackingList(pat TrackingItemQuery) (res TrackingListResponse, err error) {
	list := make([]cmn.TrackingItem, 0)

	result := s.DB.Where("user_id = ? AND status = ? AND type IN ?", pat.UserID, pat.Status, pat.Types).Find(&list)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: "No items in tracking list"}
		return
	}

	resolved := make([]cmn.MediaItem, 0, len(list))

	for _, item := range list {
		var resItem cmn.MediaItem
		resItem, err = s.ResolveItem(item.Type, item.ExternalID)

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
