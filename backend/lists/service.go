package lists

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
	"github.com/DCIAL42/lists/social/like"
	"github.com/DCIAL42/lists/users"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO: Improve error handling around all db calls

type Service struct {
	*db.DBService
	likeService *like.Service
}

func NewService(DB *gorm.DB, likeService *like.Service, clients map[cmn.MediaType]cmn.Client, config ...*db.ServiceConfig) *Service {
	return &Service{
		DBService:   db.NewDBService(DB, clients, config...),
		likeService: likeService,
	}
}

func (s *Service) toListResponse(list List, userID string) (res ListResponse, err error) {
	resolvedItems := make([]cmn.MediaResponse, 0, len(list.Items))

	for _, item := range list.Items {
		var resolvedItem cmn.MediaResponse
		resolvedItem, err = s.ResolveMedia(item.MediaID)

		if err != nil {
			err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "Error with api"}
			return
		}

		if userID != "" {
			resolvedItem.Tracking, err = s.GetTrackingItem(cmn.TrackingItem{UserID: userID, MediaID: resolvedItem.ID})
			resolvedItem.Rating, err = s.GetRating(cmn.Rating{UserID: userID, MediaID: resolvedItem.ID})
		}

		resolvedItems = append(resolvedItems, resolvedItem)
	}

	userDetails, err := users.GetUserDetails(userID)

	if err != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: "Failed to fetch user details"}
		return
	}
	likeCount := s.DB.Model(list).Association("Likes").Count()

	return ListResponse{
		ListMeta: ListMeta{
			ID:        list.ID,
			Title:     list.Title,
			CreatedBy: *userDetails.Username,
			Cover:     list.Cover,
			Likes:     uint(likeCount),
		},
		Items: resolvedItems,
	}, nil
}

func (s *Service) createList(req ListRequest, userID string) (res ListResponse, err error) {
	var x cmn.Media
	s.DB.Where("id = ?", req.Items[0].MediaID).First(&x)

	// var first cmn.MediaResponse
	// first, err = s.ResolveMedia(req.Items[0].Type, req.Items[0].ExternalID)
	items := make([]ListItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, ListItem{
			MediaID: item.MediaID,
		})
	}
	list := List{
		UserID: userID,
		Title:  req.Title,
		Items:  items,
	}
	list.Cover = x.Cover
	result := s.DB.Create(&list)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	return s.toListResponse(list, userID)
}

func (s *Service) deleteList(id uint, userID string) (ListResponse, error) {
	var list List

	result := s.DB.Clauses(clause.Returning{}).Where("id = ? AND user_id = ?", id, userID).Delete(&list)

	if result.Error != nil {
		return ListResponse{}, result.Error
	}

	if result.RowsAffected == 0 {
		return ListResponse{}, &cmn.HttpError{Code: http.StatusNotFound, Message: "not found"}
	}

	return s.toListResponse(list, list.UserID)
}

func (s *Service) updateList(id uint, userID string, req UpdateListRequest) (res ListResponse, err error) {
	var list List
	result := s.DB.
		Model(&List{}).
		Where("id = ? AND user_id = ?", id, userID).
		Preload("Items").
		First(&list)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ListResponse{}, &cmn.HttpError{
				Code:    http.StatusNotFound,
				Message: "not found",
			}
		}

		return ListResponse{}, &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
	}

	first, _ := s.ResolveItem(req.Items[0].Media.Type, req.Items[0].Media.ExternalID)

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		if req.Title != "" {
			list.Title = req.Title
		}

		if len(req.Items) != 0 {
			list.Cover = first.Cover

			if err := tx.Unscoped().
				Where("list_id = ?", list.ID).
				Delete(&ListItem{}).Error; err != nil {
				return err
			}

			list.Items = req.Items
		}

		return tx.Save(&list).Error
	})

	if err != nil {
		return res, &cmn.HttpError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	s.DB.Preload("Items").First(&list, id)

	return s.toListResponse(list, list.UserID)
}

func (s *Service) getListById(id uint, userID string) (res ListResponse, err error) {
	var list List

	result := s.DB.Preload("Items").First(&list, id)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: fmt.Sprintf("No list with id: %d", id)}
		return
	}

	return s.toListResponse(list, userID)
}

type Settings struct {
	Page  uint
	Limit uint
	Full  bool
}

func (s *Service) getListsPreviewByUser(userID string, cfg *Settings) (res ListsPreviewResponse, err error) {
	if cfg.Page < 1 {
		cfg.Page = 1
	}
	res.Page = cfg.Page
	lists := make([]List, 0)

	var result *gorm.DB
	result, res.Count = s.GetPage(cfg.Page, standardOrder("id", "desc"), &List{})
	if cfg.Limit > 0 {
		result = result.Limit(int(cfg.Limit))
	}
	result = result.Where("user_id = ?", userID).Find(&lists)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: fmt.Sprintf("No list for user: %s", userID)}
		return
	}

	res.Lists = make([]ListMeta, 0, len(lists))

	userDetails, err := users.GetUserDetails(userID)

	for _, list := range lists {
		likeCount := s.DB.Model(list).Association("Likes").Count()
		listPreview := ListMeta{
			list.ID,
			list.Title,
			*userDetails.Username,
			list.Cover,
			uint(likeCount),
		}

		res.Lists = append(res.Lists, listPreview)
	}

	return
}

func (s *Service) getListsByUser(userID string, page uint) (res ListsResponse, err error) {
	res.Page = page
	lists := make([]List, 0)

	var result *gorm.DB
	result, res.Count = s.GetPage(page, standardOrder("id", "desc"), &List{})
	result = result.Where("user_id = ?", userID).Preload("Items").Find(&lists)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: fmt.Sprintf("No list for user: %s", userID)}
		return
	}

	res.Lists = make([]ListResponse, 0, len(lists))

	for _, list := range lists {
		listResponse, err := s.toListResponse(list, list.UserID)
		likeCount := s.DB.Model(list).Association("Likes").Count()
		listResponse.Likes = uint(likeCount)

		if err != nil {
			slog.Error("Failed converting to list response", "error:", err.Error())
			continue
		}

		res.Lists = append(res.Lists, listResponse)
	}

	return
}

type ListQueryCfg struct {
	Order string
	By    string
}

func orderByLikes(order string) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Select("lists.*, COUNT(likes.id) AS like_count").
			Joins("LEFT JOIN likes ON likes.list_id = lists.id AND likes.deleted_at IS NULL").
			Group("lists.id").
			Order("like_count " + order)
	}
}

func standardOrder(by, order string) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		return d.Order(by + " " + order)
	}
}

func (s *Service) getAllLists(page uint, cfg ListQueryCfg) (res ListsPreviewResponse, err error) {
	res.Page = page
	lists := make([]List, 0)

	if cfg.Order == "" {
		cfg.Order = "desc"
	} else if cfg.Order != "asc" && cfg.Order != "desc" {
		return ListsPreviewResponse{}, &cmn.HttpError{Code: http.StatusBadRequest, Message: "order should be asc or desc"}
	}

	if cfg.By == "" {
		cfg.By = "id"
	}

	var result *gorm.DB
	if cfg.By == "likes" {
		result, res.Count = s.GetPage(page, orderByLikes(cfg.Order), &List{})
	} else {
		result, res.Count = s.GetPage(page, standardOrder(cfg.By, cfg.Order), &List{})
	}
	result = result.Preload("Likes").Find(&lists)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	res.Lists = make([]ListMeta, 0, len(lists))

	for _, list := range lists {
		userDetails, err := users.GetUserDetails(list.UserID)
		if err != nil {
			return ListsPreviewResponse{}, err
		}
		likeCount := s.DB.Model(list).Association("Likes").Count()
		res.Lists = append(res.Lists, ListMeta{
			ID:        list.ID,
			Title:     list.Title,
			CreatedBy: *userDetails.Username,
			Cover:     list.Cover,
			Likes:     uint(likeCount),
		})
	}

	return
}
