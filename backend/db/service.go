package db

import (
	"errors"

	"github.com/DCIAL42/media/cmn"
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
