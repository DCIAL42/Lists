package lists

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NextPage(u *url.URL) string {
	q := u.Query()
	page := q.Get("page")
	page_num, err := strconv.Atoi(page)
	if err != nil {
		page_num = 1
	}
	page_num++
	q.Set("page", strconv.Itoa(page_num))
	u.RawQuery = q.Encode()
	return u.String()
}

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

		if err != nil {
			// var httpErr *HttpError
			// if errors.As(err, &httpErr) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			// }
			return
		}

		c.IndentedJSON(http.StatusOK, list)
	})

	r.GET("/", func(c *gin.Context) {
		pageStr, ok := c.GetQuery("page")

		var page uint

		if ok {
			val, err := strconv.ParseUint(pageStr, 10, 64)

			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
				return
			}

			page = uint(val)
		} else {
			page = 1
		}

		lists, err := s.getAllLists(uint(page))

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		next := NextPage(c.Request.URL)
		c.IndentedJSON(http.StatusOK, gin.H{"lists": lists, "next": next})
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
