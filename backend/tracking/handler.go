package tracking

import (
	"fmt"
	"net/http"

	"github.com/DCIAL42/media/cmn"
	"github.com/DCIAL42/media/middleware"
	"github.com/gin-gonic/gin"
)

var validStatus = map[string]bool{
	"done":    true,
	"paused":  true,
	"backlog": true,
}

func (s *Service) SetupRoutes(r *gin.RouterGroup) {
	protected := r.Group("/")
	protected.Use(middleware.RequireUser())

	protected.POST("/", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)

		var item cmn.TrackingItem

		if err := c.ShouldBindJSON(&item); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item.UserID = userID

		res, err := s.createTrackingItem(item)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})

	protected.GET("/:id", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)

		id, err := cmn.ParseParam[uint](c, "id")

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var item cmn.TrackingItem

		if err := c.ShouldBindJSON(&item); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item.ID = id

		item.UserID = userID

		res, err := s.GetTrackingItem(item)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})

	protected.PATCH("/:id", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)

		id, err := cmn.ParseParam[uint](c, "id")

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var item cmn.TrackingItem

		if err := c.ShouldBindJSON(&item); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		item.ID = id

		item.UserID = userID

		res, err := s.updateTrackingItem(item)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})

	protected.DELETE("/:id", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)

		id, err := cmn.ParseParam[uint](c, "id")

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		res, err := s.deleteTrackingItem(id, userID)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})

	protected.GET("/", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)

		mediaType := c.Query("type")

		status := c.Query("status")

		if !validStatus[status] {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}

		if mediaType == "" || status == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid query params"})
			return
		}

		pat := cmn.TrackingItem{
			UserID: userID,
			Status: cmn.TrackingStatus(status),
			Type:   cmn.MediaType(mediaType),
		}
		fmt.Printf("%+v\n", pat)

		list, err := s.getTrackingList(pat)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}
		fmt.Println(list)

		c.IndentedJSON(http.StatusOK, list)
	})

	r.GET("/dev", func(c *gin.Context) {
		items, err := s.getAllTrackingItems()

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, items)
	})

	r.GET("/dev/:type/:status", func(c *gin.Context) {
		mediaType := c.Param("type")

		status := c.Param("status")

		if !validStatus[status] {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}

		pat := cmn.TrackingItem{
			Status: cmn.TrackingStatus(status),
			Type:   cmn.MediaType(mediaType),
		}

		list, err := s.getTrackingList(pat)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, list)
	})
}
