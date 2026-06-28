package search

import (
	"fmt"
	"net/http"

	"github.com/DCIAL42/media/internals/client"
	"github.com/gin-gonic/gin"
)

func NewSearchService(clients ...*client.Client) SearchService {
	return SearchService{clients}
}

func (s *SearchService) Search(c *gin.Context) {
	results := make([]client.SearchResult, 0)

	var queryParams QueryParams
	test := c.DefaultQuery("query", "empty")
	fmt.Println(test)

	err := c.BindQuery(&queryParams)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad query"})
		return
	}

	fmt.Printf("%+v\n", queryParams)

	for _, client := range s.clients {
		r, err := client.Search(c.Request.Context(), queryParams.Query)

		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "upstream server error"})
			return
		}

		results = append(results, r...)
	}

	c.IndentedJSON(http.StatusOK, results)
}
