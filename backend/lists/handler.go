package lists

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Service) SetupRoutes(r *gin.RouterGroup) {
	r.GET("/:id", func(c *gin.Context) {
		idStr := c.Param("id")

		val, err := strconv.ParseUint(idStr, 10, 64)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		id := uint(val)

		list, err := s.getListById(id)

		fmt.Println(list, err)

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		c.IndentedJSON(http.StatusOK, list)
	})

	r.GET("/", func(c *gin.Context) {
		lists, err := s.getAllLists()

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(http.StatusOK, lists)
	})

	r.POST("/", func(c *gin.Context) {
		var body List

		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		list, err := s.createList(body)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(http.StatusOK, list)
	})
}
