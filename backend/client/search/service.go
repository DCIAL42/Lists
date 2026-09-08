package search

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

func NewService(clients map[cmn.MediaType]cmn.Client, DB *gorm.DB) Service {
	return Service{
		DBService: db.NewDBService(DB, clients),
	}
}

func (s *Service) Search(c *gin.Context) {
	userID := c.GetString("userID")

	var queryParams QueryParams

	err := c.BindQuery(&queryParams)

	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad query"})
		return
	}

	var resultTypes []string
	if queryParams.Types == "all" {
		resultTypes = []string{"movie", "album"}
	} else {
		resultTypes = strings.Split(string(queryParams.Types), "|")
	}

	g, ctx := errgroup.WithContext(c.Request.Context())
	ctx = context.WithValue(ctx, "originalURL", c.Request.RequestURI)
	ctx = context.WithValue(ctx, "userID", userID)

	results := make(map[cmn.MediaType]cmn.SearchResult)

	for _, resultType := range resultTypes {
		cl := s.Clients[cmn.MediaType(resultType)]

		g.Go(func() error {
			r, err := cl.Search(ctx, map[string]string{"query": queryParams.Query, "page": queryParams.Page})

			if err != nil {
				slog.Error(err.Error())
				return nil
			}

			results[cmn.MediaType(resultType)] = r

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		slog.Error(err.Error())
	}

	c.IndentedJSON(http.StatusOK, results)
}
