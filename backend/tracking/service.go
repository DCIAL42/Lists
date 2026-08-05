package tracking

import (
	"net/http"

	"github.com/DCIAL42/media/cmn"
	"github.com/DCIAL42/media/db"
	"github.com/DCIAL42/media/internals/client"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	*db.DBService
}

func NewService(DB *gorm.DB, clients map[cmn.MediaType]client.Client, config ...*db.ServiceConfig) *Service {
	return &Service{
		DBService: db.NewDBService(DB, clients, config...),
	}
}

func (s *Service) createTrackingItem(req TrackingItem) (res TrackingItem, err error) {
	result := s.DB.Create(&req)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	res = req

	return
}

func (s *Service) updateTrackingItem(req TrackingItem) (res TrackingItem, err error) {
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

func (s *Service) deleteTrackingItem(id uint, userID string) (res TrackingItem, err error) {
	result := s.DB.Clauses(clause.Returning{}).Unscoped().Where("id = ? AND user_id = ?", id, userID).Delete(&res)

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

func (s *Service) getTrackingList(pat TrackingItem) (res TrackingListResponse, err error) {
	list := make([]TrackingItem, 0)

	result := s.DB.Where(&pat).Find(&list)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: "No items in tracking list"}
		return
	}

	resolved := make([]TrackingItemResponse, 0, len(list))

	for _, item := range list {
		var resItem cmn.MediaItem
		resItem, err = s.ResolveItem(pat.Type, item.ExternalID)

		if err != nil {
			err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "Error with api"}
			return
		}

		resolved = append(resolved, TrackingItemResponse{
			ID:        item.ID,
			MediaItem: resItem,
		})
	}

	res = TrackingListResponse{
		Items: resolved,
	}

	return
}

func (s *Service) getAllTrackingItems() (res []TrackingItem, err error) {
	res = make([]TrackingItem, 0)

	result := s.DB.Find(&res)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: "No items in tracking list"}
		return
	}

	return
}
