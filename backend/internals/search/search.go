package search

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/DCIAL42/media/internals/client"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func NewSearchService(clients ...*client.Client) SearchService {
	return SearchService{clients}
}

func (s *SearchService) Search(c *gin.Context) {
	results := make([]client.SearchResult, 0)

	var queryParams QueryParams

	err := c.BindQuery(&queryParams)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad query"})
		return
	}

	resultTypes := strings.Split(queryParams.Types, "|")

	var mu sync.Mutex
	g, ctx := errgroup.WithContext(c.Request.Context())

	for _, client := range s.clients {
		if !slices.Contains(resultTypes, client.GetType()) {
			continue
		}

		g.Go(func() error {
			r, err := client.Search(ctx, map[string]string{"query": queryParams.Query})

			if err != nil {
				return err
			}

			mu.Lock()
			results = append(results, r...)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}

	c.IndentedJSON(http.StatusOK, results)
}
