package search

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/DCIAL42/lists/cmn"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

func NewSearchService(clients map[cmn.MediaType]cmn.Client, db *gorm.DB) Service {
	return Service{clients, db}
}

func AddTrackingInfo(db *gorm.DB, items []cmn.MediaItem, userID string) (res []cmn.MediaItem, err error) {
	res = make([]cmn.MediaItem, 0, len(items))

	for _, item := range items {
		var tracking cmn.TrackingItem
		db.Where("external_id = ? AND user_id = ?", item.ExternalID, userID).First(&tracking)

		item.Tracking = cmn.TrackingResponse{
			ID:     tracking.ID,
			Status: tracking.Status,
		}

		res = append(res, item)
	}

	return
}

func (s *Service) Search(c *gin.Context) {
	claims, ok := clerk.SessionClaimsFromContext(c.Request.Context())

	results := make([]cmn.SearchResult, 0)

	var queryParams QueryParams

	err := c.BindQuery(&queryParams)

	if err != nil {
		slog.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad query"})
		return
	}

	resultTypes := strings.Split(string(queryParams.Types), "|")

	var mu sync.Mutex
	g, ctx := errgroup.WithContext(c.Request.Context())
	ctx = context.WithValue(ctx, "originalURL", c.Request.RequestURI)

	for _, resultType := range resultTypes {
		cl := s.clients[cmn.MediaType(resultType)]

		g.Go(func() error {
			r, err := cl.Search(ctx, map[string]string{"query": queryParams.Query, "page": queryParams.Page})

			if err != nil {
				slog.Error(err.Error())
				return err
			}

			if ok && claims.Subject != "" {
				r.Items, err = AddTrackingInfo(s.db, r.Items, claims.Subject)
			}

			mu.Lock()
			results = append(results, r)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		slog.Error(err.Error())
	}

	if len(results) == 1 {
		c.IndentedJSON(http.StatusOK, results[0])
	} else {
		c.IndentedJSON(http.StatusOK, results)
	}
}
