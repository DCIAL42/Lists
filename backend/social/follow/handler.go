package follow

import (
	"net/http"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/middleware"
	"github.com/DCIAL42/lists/users"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	*gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) SetupRoutes(r *gin.RouterGroup) {
	protected := r.Group("/")
	protected.Use(middleware.RequireUser())

	protected.DELETE("/:id", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)
		id, err := cmn.ParseParam[uint](c, "id")

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		res, err := s.deleteFollow(userID, id)

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

	protected.POST("/follow", func(c *gin.Context) {
		follower := c.MustGet("userID").(string)

		username := c.Param("username")
		followed, err := users.GetUserByUsername(username)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		res, err := s.newFollow(follower, followed.ID)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})

	protected.GET("/follow", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)

		username := c.Param("username")
		followed, err := users.GetUserByUsername(username)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		res, err := s.getFollow(userID, followed.ID)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}
		c.IndentedJSON(http.StatusOK, res)
	})

	r.GET("/following", func(c *gin.Context) {
		username := c.Param("username")
		user, err := users.GetUserByUsername(username)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		res, err := s.getFollowers(user.ID)
		if err != nil {
			cmn.HandleError(c, err)
			return
		}
		c.IndentedJSON(http.StatusOK, res)
	})
}
