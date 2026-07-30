package lists

import (
	"errors"
	"fmt"

	"github.com/DCIAL42/media/internals/music"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO: Improve error handling around all db calls

type Service struct {
	db          *gorm.DB
	musicClient *music.Client
}

func NewService(musicClient *music.Client) *Service {
	return &Service{
		musicClient: musicClient,
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

func (s *Service) resolveItem(item ListItem) (any, error) {
	switch item.Type {
	case ItemTypeAlbum:
		return s.musicClient.GetAlbum(item.ExternalID)
	}
	return nil, errors.New("Invalid item type, unable to resolve.")
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
		return ListResponse{}, err
	}

	resolved := make([]any, 0)

	for _, item := range list.Items {
		res, err := s.resolveItem(item)

		if err != nil {
			continue
		}

		fmt.Printf("%+v", res)
		resolved = append(resolved, res)
	}

	res := ListResponse{
		Title:     list.Title,
		CreatedBy: list.CreatedBy,
		Items:     resolved,
	}

	return res, nil
}

func (s *Service) getAllLists() ([]List, error) {
	db, err := initDB()

	if err != nil {
		return []List{}, err
	}

	lists := make([]List, 0)

	result := db.Find(&lists)

	if result.Error != nil {
		return []List{}, err
	}

	fmt.Println(lists)

	return lists, nil

	// resolved := make([]any, 0)
	// for _, item := range list.Items {
	// 	res, err := s.resolveItem(item)
	//
	// 	if err != nil {
	// 		continue
	// 	}
	//
	// 	fmt.Printf("%+v", res)
	// 	resolved = append(resolved, res)
	// }
	//
	// res := ListResponse{
	// 	Title:     list.Title,
	// 	CreatedBy: list.CreatedBy,
	// 	Items:     resolved,
	// }
	//
	// return res, nil
}
