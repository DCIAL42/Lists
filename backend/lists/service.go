package lists

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
	"github.com/DCIAL42/lists/users"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO: Improve error handling around all db calls

type Service struct {
	*db.DBService
}

func NewService(DB *gorm.DB, clients map[cmn.MediaType]cmn.Client, config ...*db.ServiceConfig) *Service {
	return &Service{
		DBService: db.NewDBService(DB, clients, config...),
	}
}

func (s *Service) toListResponse(list List) (res ListResponse, err error) {
	resolvedItems := make([]cmn.MediaItem, 0, len(list.Items))

	for _, item := range list.Items {
		var resolvedItem cmn.MediaItem
		resolvedItem, err = s.ResolveItem(item.Type, item.ExternalID)

		if err != nil {
			err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "Error with api"}
			return
		}

		resolvedItems = append(resolvedItems, resolvedItem)
	}

	userDetails, err := users.GetUserDetails(list.UserID)

	if err != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "Failed to fetch user details"}
		return
	}

	return ListResponse{
		ID:        list.ID,
		Title:     list.Title,
		CreatedBy: *userDetails.Username,
		Items:     resolvedItems,
	}, nil
}

func (s *Service) createList(list List) (res ListResponse, err error) {
	result := s.DB.Create(&list)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	return s.toListResponse(list)
}

func (s *Service) deleteList(id uint, userID string) (ListResponse, error) {
	var list List

	result := s.DB.Clauses(clause.Returning{}).Unscoped().Where("id = ? AND user_id = ?", id, userID).Delete(&list)

	if result.Error != nil {
		return ListResponse{}, result.Error
	}

	if result.RowsAffected == 0 {
		return ListResponse{}, &cmn.HttpError{Code: http.StatusNotFound, Message: "not found"}
	}

	return s.toListResponse(list)
}

func (s *Service) getListById(id uint) (res ListResponse, err error) {
	var list List

	result := s.DB.Preload("Items").First(&list, id)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: fmt.Sprintf("No list with id: %d", id)}
		return
	}

	return s.toListResponse(list)
}

func (s *Service) getAllLists(page uint) (res []ListResponse, err error) {
	lists := make([]List, 0)

	result := s.DB.Limit(int(s.PageSize)).Offset((int(page) - 1) * int(s.PageSize)).Preload("Items").Find(&lists)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	res = make([]ListResponse, 0, len(lists))

	for _, list := range lists {
		listResponse, err := s.toListResponse(list)

		if err != nil {
			slog.Error("Failed converting to list response", "error:", err.Error())
			continue
		}

		res = append(res, listResponse)
	}

	return
}
