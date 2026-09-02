package search

import (
	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
)

type Service struct {
	*db.DBService
}

type QueryParams struct {
	Query string        `form:"query"`
	Types cmn.MediaType `form:"type"`
	Page  string        `form:"page"`
	Full  bool          `form:"full"`
}
