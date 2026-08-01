package lists

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DCIAL42/media/internals/client"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO: Improve error handling around all db calls

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

func NewService(musicClient client.Client, movieClient client.Client) *Service {
	return &Service{
		musicClient: musicClient,
		movieClient: movieClient,
	}
}

func initDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&List{},
		&ListItem{},
	)

	return db, err
}

func (s *Service) resolveItem(item ListItem) (client.MediaItem, error) {
	switch item.Type {
	case client.TypeAlbum:
		album, err := s.musicClient.GetItem(item.ExternalID)

		if err != nil {
			return client.MediaItem{}, err
		}

		return album, nil
	case client.TypeMovie:
		movie, err := s.movieClient.GetItem(item.ExternalID)

		if err != nil {
			return client.MediaItem{}, err
		}

		return movie, nil
	}
	return client.MediaItem{}, errors.New("Invalid item type, unable to resolve.")
}

func (s *Service) createList(req List) (List, error) {
	db, err := initDB()

	if err != nil {
		return List{}, err
	}

	result := db.Create(&req)

	if result.Error != nil {
		return List{}, result.Error
	}

	return req, nil
}

// TODO: Add resolving of list items into full data
func (s *Service) getListById(id uint) (ListResponse, error) {
	db, err := initDB()

	if err != nil {
		return ListResponse{}, err
	}

	var list List

	result := db.Preload("Items").First(&list, id)

	if result.Error != nil {
		return ListResponse{}, &HttpError{Code: http.StatusNotFound, Message: fmt.Sprintf("No list with id: %d", id)}
	}

	resolved := make([]client.MediaItem, 0)

	for _, item := range list.Items {
		res, err := s.resolveItem(item)

		if err != nil {
			return ListResponse{}, &HttpError{Code: http.StatusInternalServerError, Message: "Error with api"}
		}

		resolved = append(resolved, res)
	}

	res := ListResponse{
		Title:     list.Title,
		CreatedBy: list.CreatedBy,
		Items:     resolved,
	}

	return res, nil
}

func (s *Service) getAllLists(page uint) ([]ListResponse, error) {
	db, err := initDB()

	if err != nil {
		return []ListResponse{}, err
	}

	lists := make([]List, 0)

	result := db.Limit(10).Offset((int(page) - 1) * 10).Preload("Items").Find(&lists)

	if result.Error != nil {
		return []ListResponse{}, err
	}

	res := make([]ListResponse, 0, len(lists))

	for _, list := range lists {
		items := make([]client.MediaItem, 0, len(list.Items))

		for _, item := range list.Items {
			resolved, err := s.resolveItem(item)

			if err != nil {
				continue
				// return []ListResponse{}, err
			}

			items = append(items, client.MediaItem{
				Type: item.Type,
				Data: resolved,
			})
		}

		res = append(res, ListResponse{
			Title:     list.Title,
			CreatedBy: list.CreatedBy,
			Items:     items,
		})
	}

	return res, nil
}
