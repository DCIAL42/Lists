package search

import (
	"github.com/DCIAL42/lists/cmn"
	"gorm.io/gorm"
)

type Service struct {
	clients map[cmn.MediaType]cmn.Client
	db      *gorm.DB
}

type QueryParams struct {
	Query string        `form:"query"`
	Types cmn.MediaType `form:"type"`
	Page  string        `form:"page"`
	Full  bool          `form:"full"`
}
