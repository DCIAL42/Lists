package lists

import (
	"net/http"
	"strconv"

	"github.com/DCIAL42/media/cmn"
	"github.com/DCIAL42/media/middleware"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gin-gonic/gin"
)

func (s *Service) SetupRoutes(r *gin.RouterGroup) {
	protected := r.Group("/")
	protected.Use(
		middleware.WrapClerkMiddleware(clerkhttp.WithHeaderAuthorization()),
		middleware.RequireUser(),
	)

	r.GET("/:id", func(c *gin.Context) {
		id, err := cmn.ParseParam[uint](c, "id")

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

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

		next := cmn.NextPage(c.Request.URL)
		c.IndentedJSON(http.StatusOK, gin.H{"lists": lists, "next": next})
	})

	protected.POST("/", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)

		var body List

		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		body.UserID = userID

		list, err := s.createList(body)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(http.StatusOK, list)
	})
}
