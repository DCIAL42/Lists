package db

import (
	"errors"
	"net/http"

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

func (db *DBService) GetPage(page uint, order string, dst any) (*gorm.DB, uint) {
	var count int64
	result := db.DB.
		Model(dst).
		Order(order).
		Count(&count).
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

type ExternalItem interface {
	GetExternalID() string
}

func TrySaveItem[T ExternalItem](DB *gorm.DB, dst *T) (bool, error) {
	var existing T
	result := DB.Where("external_id = ?", (*dst).GetExternalID()).First(&existing)

	if result.Error == nil {
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
	result := DB.Where("external_id = ?", externalID).First(dst)
	if result.Error != nil {
		return false
	}
	return true
}
