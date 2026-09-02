package search

import (
	"github.com/DCIAL42/lists/middleware"
	"github.com/gin-gonic/gin"
)

func (s *Service) SetupRoutes(r *gin.RouterGroup) {
	maybe := r.Group("/")
	maybe.Use(middleware.MaybeUser())

	maybe.GET("/", s.Search)
}
