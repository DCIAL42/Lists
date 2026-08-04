package tracking

import (
	"errors"
	"net/http"

	"github.com/DCIAL42/media/cmn"
	"github.com/DCIAL42/media/internals/client"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HttpError struct {
	Code    int
	Message string
	Err     error
}

func (e *HttpError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

type Service struct {
	db          *gorm.DB
	musicClient client.Client
	movieClient client.Client
}

func NewService(db *gorm.DB, musicClient client.Client, movieClient client.Client) *Service {
	return &Service{
		db:          db,
		musicClient: musicClient,
		movieClient: movieClient,
	}
}

func (s *Service) resolveItem(t cmn.MediaType, externalID string) (cmn.MediaItem, error) {
	switch t {
	case cmn.TypeAlbum:
		album, err := s.musicClient.GetItem(externalID)

		if err != nil {
			return cmn.MediaItem{}, err
		}

		return album, nil
	case cmn.TypeMovie:
		movie, err := s.movieClient.GetItem(externalID)

		if err != nil {
			return cmn.MediaItem{}, err
		}

		return movie, nil
	}
	return cmn.MediaItem{}, errors.New("Invalid item type, unable to resolve.")
}

func (s *Service) createTrackingItem(req TrackingItem) (TrackingItem, error) {
	result := s.db.Create(&req)

	if result.Error != nil {
		return TrackingItem{}, result.Error
	}

	return req, nil
}

func (s *Service) updateTrackingItem(req TrackingItem) (TrackingItem, error) {
	result := s.db.Save(&req)

	if result.Error != nil {
		return TrackingItem{}, result.Error
	}

	return req, nil
}

func (s *Service) deleteTrackingItem(id uint, userID string) (TrackingItem, error) {
	var item TrackingItem

	result := s.db.Clauses(clause.Returning{}).Unscoped().Where("id = ? AND user_id = ?", id, userID).Delete(&item)

	if result.Error != nil {
		return TrackingItem{}, result.Error
	}

	if result.RowsAffected == 0 {
		return TrackingItem{}, errors.New("Item not found")
	}

	return item, nil
}

func (s *Service) getTrackingList(pat TrackingItem) (TrackingListResponse, error) {
	list := make([]TrackingItem, 0)

	result := s.db.Where(&pat).Find(&list)

	if result.Error != nil {
		return TrackingListResponse{}, &HttpError{Code: http.StatusNotFound, Message: "No items in tracking list"}
	}

	resolved := make([]TrackingItemResponse, 0, len(list))

	for _, item := range list {
		res, err := s.resolveItem(pat.Type, item.ExternalID)

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
