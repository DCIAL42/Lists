package users

import (
	"net/http"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	*gorm.DB
}

func NewUserService(db *gorm.DB) *Service {
	return &Service{db}
}

func (s *Service) SetupRoutes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		users, err := s.getAllUsers()
		if err != nil {
			cmn.HandleError(c, err)
			return
		}
		c.IndentedJSON(http.StatusOK, users)
	})
}

func (s *Service) SetupUserRoutes(r *gin.RouterGroup) {
	protected := r.Group("/")
	protected.Use(middleware.RequireUser())

	protected.POST("/follow", func(c *gin.Context) {
		fromID := c.MustGet("userID").(string)

		username := c.Param("username")
		to, err := GetUserByUsername(username)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}

		res, err := s.newFollowing(fromID, to.ID)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})

	r.GET("/", func(c *gin.Context) {
		username := c.Param("username")
		user, err := GetUserByUsername(username)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.IndentedJSON(http.StatusOK, UserResponse{ID: user.ID, Username: *user.Username})
	})

	r.GET("/following", func(c *gin.Context) {
		username := c.Param("username")
		user, err := GetUserByUsername(username)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		res, err := s.getFollowing(user.ID)
		if err != nil {
			cmn.HandleError(c, err)
			return
		}
		c.IndentedJSON(http.StatusOK, res)
	})
}
