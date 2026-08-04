package tracking

import (
	"net/http"

	"github.com/DCIAL42/media/cmn"
	"github.com/DCIAL42/media/db"
	"github.com/DCIAL42/media/internals/client"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HttpError struct {
	Code    int
	Message string
}

func (e *HttpError) Error() string {
	return e.Message
}

type Service struct {
	*db.DBService
}

func NewService(DB *gorm.DB, clients map[cmn.MediaType]client.Client, config ...*db.ServiceConfig) *Service {
	return &Service{
		DBService: db.NewDBService(DB, clients, config...),
	}
}

func (s *Service) createTrackingItem(req TrackingItem) (TrackingItem, error) {
	result := s.DB.Create(&req)

	if result.Error != nil {
		return TrackingItem{}, result.Error
	}

	return req, nil
}

func (s *Service) updateTrackingItem(req TrackingItem) (TrackingItem, error) {
	result := s.DB.Where("id = ? AND user_id = ?", req.ID, req.UserID).Updates(req)

	if result.Error != nil {
		return TrackingItem{}, result.Error
	}

	if result.RowsAffected == 0 {
		return TrackingItem{}, &cmn.HttpError{Code: http.StatusNotFound, Message: "not found"}
	}

	var res TrackingItem
	s.DB.First(&res, req.ID)

	return res, nil
}

func (s *Service) deleteTrackingItem(id uint, userID string) (TrackingItem, error) {
	var item TrackingItem

	result := s.DB.Clauses(clause.Returning{}).Unscoped().Where("id = ? AND user_id = ?", id, userID).Delete(&item)

	if result.Error != nil {
		return TrackingItem{}, result.Error
	}

	if result.RowsAffected == 0 {
		return TrackingItem{}, &cmn.HttpError{Code: http.StatusNotFound, Message: "not found"}
	}

	return item, nil
}

func (s *Service) getTrackingList(pat TrackingItem) (TrackingListResponse, error) {
	list := make([]TrackingItem, 0)

	result := s.DB.Where(&pat).Find(&list)

	if result.Error != nil {
		return TrackingListResponse{}, &HttpError{Code: http.StatusNotFound, Message: "No items in tracking list"}
	}

	resolved := make([]TrackingItemResponse, 0, len(list))

	for _, item := range list {
		res, err := s.ResolveItem(pat.Type, item.ExternalID)

		if err != nil {
			return TrackingListResponse{}, &HttpError{Code: http.StatusInternalServerError, Message: "Error with api"}
		}

		resolved = append(resolved, TrackingItemResponse{
			ID:        item.ID,
			MediaItem: res,
		})
	}

	res := TrackingListResponse{
		Items: resolved,
	}

	return res, nil
}

func (s *Service) getAllTrackingItems() ([]TrackingItem, error) {
	list := make([]TrackingItem, 0)

	result := s.DB.Find(&list)

	if result.Error != nil {
		return []TrackingItem{}, &HttpError{Code: http.StatusNotFound, Message: "No items in tracking list"}
	}

	return list, nil
}
