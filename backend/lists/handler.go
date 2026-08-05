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
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, list)
	})

	r.GET("/", func(c *gin.Context) {
		pageStr := c.DefaultQuery("page", "1")

		val, err := strconv.ParseUint(pageStr, 10, 64)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
			return
		}

		page := uint(val)

		lists, err := s.getAllLists(page)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		next := cmn.NextPage(c.Request.URL)
		c.IndentedJSON(http.StatusOK, gin.H{"lists": lists, "next": next})
	})

	protected.POST("/", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)

		var body List

		if err := c.ShouldBindJSON(&body); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		body.UserID = userID

		list, err := s.createList(body)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, list)
	})

	protected.DELETE("/:id", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)

		id, err := cmn.ParseParam[uint](c, "id")

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		s.deleteList(id, userID)
	})
}
