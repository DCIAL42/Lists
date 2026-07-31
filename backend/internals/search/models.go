package search

import "github.com/DCIAL42/media/internals/client"

type SearchService struct {
	clients []client.Client
}

type QueryParams struct {
	Query string           `form:"query"`
	Types client.MediaType `form:"type"`
	Page  string           `form:"page"`
}
