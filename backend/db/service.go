package db

import (
	"errors"
	"net/http"
	"time"

	"github.com/DCIAL42/lists/cmn"
	"gorm.io/gorm"
)

type ServiceConfig struct {
	PageSize uint
}

type DBService struct {
	*ServiceConfig
	DB      *gorm.DB
	Clients map[cmn.MediaType]cmn.Client
}

func defaultDBServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		PageSize: 10,
	}
}

func (s *DBService) ResolveItem(t cmn.MediaType, externalID string) (cmn.MediaItem, error) {
	c, ok := s.Clients[t]

	if !ok {
		return cmn.MediaItem{}, errors.New("Invalid item type, unable to resolve.")
	}

	return c.GetItem(externalID)
}

func (s *DBService) ResolveMedia(mediaID uint) (cmn.MediaResponse, error) {
	var media cmn.Media
	if err := s.DB.Where("id = ?", mediaID).First(&media).Error; err != nil {
		return cmn.MediaResponse{}, &cmn.HttpError{Code: http.StatusNotFound, Message: "media item not found"}
	}
	t := media.Type
	c, ok := s.Clients[t]

	if !ok {
		return cmn.MediaResponse{}, errors.New("Invalid item type, unable to resolve.")
	}

	return c.GetMedia(mediaID)
}

func NewDBService(db *gorm.DB, clients map[cmn.MediaType]cmn.Client, config ...*ServiceConfig) *DBService {
	var cfg *ServiceConfig

	if len(config) == 0 {
		cfg = defaultDBServiceConfig()
	} else {
		cfg = config[0]
	}

	return &DBService{
		ServiceConfig: cfg,
		DB:            db,
		Clients:       clients,
	}
}

func (db *DBService) GetPage(page uint, order func(*gorm.DB) *gorm.DB, dst any) (*gorm.DB, uint) {
	var count int64
	if result := db.DB.Model(dst).Count(&count); result.Error != nil {
		return result, 0
	}

	result := order(db.DB).
		Offset((int(page) - 1) * int(db.PageSize)).
		Limit(int(db.PageSize))

	return result, uint(count)
}

func (db *DBService) GetTrackingItem(req cmn.TrackingItem) (res cmn.TrackingResponse, err error) {
	var item cmn.TrackingItem
	result := db.DB.Where(req).First(&item)

	if result.Error != nil {
		err = &cmn.HttpError{Code: http.StatusInternalServerError, Message: result.Error.Error()}
		return
	}

	return cmn.TrackingResponse{
		ID:     item.ID,
		Status: item.Status,
	}, nil
}

func (db *DBService) GetRating(req cmn.Rating) (res cmn.RatingResponse, err error) {
	var r cmn.Rating
	err = db.DB.Where(req).First(&r).Error

	res = cmn.RatingResponse{
		ID:     r.ID,
		Rating: r.Rating,
	}

	return
}

type ExternalItem interface {
	GetExternalID() string
	GetModel() cmn.Model
}

func TrySaveItem[T ExternalItem](DB *gorm.DB, dst T) (bool, error) {
	var media cmn.Media

	mediaResult := DB.Where("external_id = ?", dst.GetExternalID()).First(&media)

	if mediaResult.Error != nil && !errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		return false, mediaResult.Error
	}

	if mediaResult.Error != nil && errors.Is(mediaResult.Error, gorm.ErrRecordNotFound) {
		if err := DB.Create(dst).Error; err != nil {
			return false, err
		}
		return true, nil
	}

	var existing T
	result := DB.Where("media_id = ?", media.ID).Preload("Media").First(&existing)

	if result.Error == nil {
		if time.Since(existing.GetModel().UpdatedAt) > time.Second {
			if err := DB.Model(&existing).Updates(dst).Error; err != nil {
				return false, err
			}

			if err := DB.Preload("Media").First(dst, existing.GetModel().ID).Error; err != nil {
				return false, err
			}
			return true, nil
		}
		return false, nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, result.Error
	}

	if err := DB.Create(dst).Error; err != nil {
		return false, err
	}

	return true, nil
}

func TryGetItem(DB *gorm.DB, externalID string, dst any) bool {
	result := DB.Where("external_id = ?", externalID).Preload("Media").First(dst)
	if result.Error != nil {
		return false
	}
	return true
}
