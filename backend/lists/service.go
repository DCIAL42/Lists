package lists

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DCIAL42/media/cmn"
	"github.com/DCIAL42/media/internals/client"
	"github.com/DCIAL42/media/users"
	"gorm.io/gorm"
)

// TODO: Improve error handling around all db calls

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

func (s *Service) toListResponse(list List) (res ListResponse, err error) {
	resolvedItems := make([]cmn.MediaItem, 0, len(list.Items))

	for _, item := range list.Items {
		resolvedItem, err := s.resolveItem(item.Type, item.ExternalID)

		if err != nil {
			return ListResponse{}, &cmn.HttpError{Code: http.StatusInternalServerError, Message: "Error with api"}
		}

		resolvedItems = append(resolvedItems, resolvedItem)
	}

	userDetails, err := users.GetUserDetails(list.UserID)

	if err != nil {
		return ListResponse{}, err
	}

	return ListResponse{
		ID:        list.ID,
		Title:     list.Title,
		CreatedBy: *userDetails.Username,
		Items:     resolvedItems,
	}, nil
}

func (s *Service) resolveItem(t cmn.MediaType, externalID string) (cmn.MediaItem, error) {
	switch t {
	case cmn.TypeAlbum:
		return s.musicClient.GetItem(externalID)
	case cmn.TypeMovie:
		return s.movieClient.GetItem(externalID)
	}
	return cmn.MediaItem{}, errors.New("Invalid item type, unable to resolve.")
}

func (s *Service) createList(list List) (ListResponse, error) {
	result := s.db.Create(&list)

	if result.Error != nil {
		return ListResponse{}, result.Error
	}

	return s.toListResponse(list)
}

func (s *Service) getListById(id uint) (res ListResponse, err error) {
	var list List

	result := s.db.Preload("Items").First(&list, id)

	if result.Error != nil {
		return ListResponse{}, &cmn.HttpError{Code: http.StatusNotFound, Message: fmt.Sprintf("No list with id: %d", id)}
	}

	return s.toListResponse(list)
}

func (s *Service) getAllLists(page uint) (res []ListResponse, err error) {
	lists := make([]List, 0)

	result := s.db.Limit(100).Offset((int(page) - 1) * 100).Preload("Items").Find(&lists)

	if result.Error != nil {
		err = result.Error
		return
	}

	res = make([]ListResponse, 0, len(lists))

	for _, list := range lists {
		listResponse, err := s.toListResponse(list)

		if err != nil {
			continue
		}

		res = append(res, listResponse)
	}

	return
}
