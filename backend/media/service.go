package media

import (
	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
	"gorm.io/gorm"
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

func (s *Service) getMedia(mediaID uint, userID string) (cmn.MediaResponse, error) {
	media, err := s.GetMedia(mediaID, userID)

	if err != nil {
		return cmn.MediaResponse{}, err
	}

	return s.ToMediaResponse(media)
}
