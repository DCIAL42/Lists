package tracking

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/middleware"
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

		var test cmn.TrackingItem

		result := s.DB.Unscoped().Where("media_id = ? AND user_id = ?", item.MediaID, userID).First(&test)

		if result.Error == nil {
			result := s.DB.Unscoped().Model(&test).Updates(map[string]any{
				"deleted_at": nil,
				"status":     item.Status,
			})
			if result.Error != nil {
				c.IndentedJSON(http.StatusInternalServerError, result.Error.Error())
				return
			}
			if result.RowsAffected == 0 {
				c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.IndentedJSON(http.StatusOK, cmn.TrackingResponse{ID: test.ID, Status: test.Status})
			return
		}

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

		res, err := s.getTrackingItem(item)

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

		mediaTypeStr := c.Query("type")

		status := c.Query("status")
		pageStr := c.DefaultQuery("page", "1")

		page, err := strconv.Atoi(pageStr)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
			return
		}

		if !validStatus[status] {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}

		if mediaTypeStr == "" || status == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid query params"})
			return
		}

		mediaTypes := make([]cmn.MediaType, 0)

		for t := range strings.SplitSeq(mediaTypeStr, "|") {
			mediaTypes = append(mediaTypes, cmn.MediaType(t))
		}

		pat := TrackingItemQuery{
			userID,
			cmn.TrackingStatus(status),
			mediaTypes,
		}

		list, err := s.getTrackingList(pat, page)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

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
}
