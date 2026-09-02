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

func (s *DBService) GetMedia(mediaID uint, userID string) (res cmn.Media, err error) {
	tx := s.DB
	if userID != "" {
		tx = tx.Preload("Tracking", "user_id = ?", userID).Preload("Rating", "user_id = ?", userID)
	}
	if err = tx.Where("id = ?", mediaID).First(&res).Error; err != nil {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: "media item not found"}
		return
	}

	return
}

func (s *DBService) GetMediaList(mediaIDs []uint, userID string) (res []cmn.Media, err error) {
	tx := s.DB
	if userID != "" {
		tx = tx.Preload("Tracking", "user_id = ?", userID).Preload("Rating", "user_id = ?", userID)
	}
	if err = tx.Where("id IN ?", mediaIDs).Find(&res).Error; err != nil {
		err = &cmn.HttpError{Code: http.StatusNotFound, Message: "media item not found"}
		return
	}

	return
}

func (s *DBService) ToMediaResponse(m cmn.Media) (res cmn.MediaResponse, err error) {
	t := m.Type
	c, ok := s.Clients[t]

	if !ok {
		return cmn.MediaResponse{}, errors.New("Invalid item type, unable to resolve.")
	}

	return c.ResolveMedia(m)
}

func (s *DBService) ResolveMedia(mediaID uint, userID string) (res cmn.MediaResponse, err error) {
	var media cmn.Media
	media, err = s.GetMedia(mediaID, userID)

	if err != nil {
		return
	}

	return s.ToMediaResponse(media)
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
		if time.Since(existing.GetModel().UpdatedAt) > time.Hour*24 {
			if err := DB.Model(&existing).Updates(dst).Error; err != nil {
				return false, err
			}

			if err := DB.Preload("Media").First(dst, existing.GetModel().ID).Error; err != nil {
				return false, err
			}
			return true, nil
		}
		if err := DB.Preload("Media").First(dst, existing.GetModel().ID).Error; err != nil {
			return false, err
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
