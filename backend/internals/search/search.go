package search

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/DCIAL42/media/internals/client"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func NewSearchService(clients ...client.Client) SearchService {
	return SearchService{clients}
}

func (s *SearchService) Search(c *gin.Context) {
	results := make([]client.SearchResult, 0)

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

	for _, client := range s.clients {
		if !slices.Contains(resultTypes, string(client.GetMediaType())) {
			continue
		}

		g.Go(func() error {
			r, err := client.Search(ctx, map[string]string{"query": queryParams.Query, "page": queryParams.Page})

			if err != nil {
				slog.Error(err.Error())
				return err
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
