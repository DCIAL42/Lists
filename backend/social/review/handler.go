package review

import (
	"net/http"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/middleware"
	"github.com/DCIAL42/lists/users"
	"github.com/gin-gonic/gin"
)

func (s *Service) SetupRoutes(r *gin.RouterGroup) {
	protected := r.Group("/")
	protected.Use(middleware.RequireUser())

	r.GET("/", func(c *gin.Context) {
		res, err := s.getAllReviews()

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})

	r.GET("/:mediaID", func(c *gin.Context) {
		mediaID, err := cmn.ParseParam[uint](c, "mediaID")

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad media id"})
			return
		}

		res, err := s.getReviewsByMedia(mediaID)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})

	r.POST("/test", func(c *gin.Context) {
		userID := "user_3HN6SMtn6hDq0JTcUlZL3Bwr7fu"

		var req ReviewRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad post data"})
		}

		res, err := s.createReview(userID, req)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})

	protected.POST("/", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)

		var req ReviewRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "bad post data"})
		}

		res, err := s.createReview(userID, req)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})
}

func (s *Service) SetupUserRoutes(r *gin.RouterGroup) {
	protected := r.Group("/")
	protected.Use(middleware.RequireUser())

	r.GET("/reviews", func(c *gin.Context) {
		username := c.Param("username")
		user, err := users.GetUserByUsername(username)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}

		res, err := s.getReviewsByUser(user.ID)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})
}
