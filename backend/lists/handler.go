package lists

import (
	"net/http"
	"strconv"

	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/middleware"
	"github.com/DCIAL42/lists/users"
	"github.com/gin-gonic/gin"
)

func (s *Service) SetupRoutes(r *gin.RouterGroup) {
	protected := r.Group("/")
	protected.Use(middleware.MaybeUser())

	protected.GET("/:id", func(c *gin.Context) {
		userID := c.GetString("userID")
		id, err := cmn.ParseParam[uint](c, "id")

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		list, err := s.getListById(id, userID)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		like, err := s.likeService.GetLike(userID, id)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"list": list, "like": like})
	})

	r.GET("/", func(c *gin.Context) {
		order := c.DefaultQuery("order", "asc")
		pageStr := c.DefaultQuery("page", "1")

		val, err := strconv.ParseUint(pageStr, 10, 64)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
			return
		}

		page := uint(val)

		reverse := true
		if order == "asc" {
			reverse = false
		}

		res, err := s.getAllLists(page, reverse)

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		res.Next = cmn.NextPage(c.Request.URL)

		c.IndentedJSON(http.StatusOK, res)
	})

	protected.POST("/", func(c *gin.Context) {
		userID := c.MustGet("userID").(string)

		var body List

		if err := c.ShouldBindJSON(&body); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if body.Title == "" {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "empty title field"})
			return
		}

		if len(body.Items) == 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no items in list"})
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

		res, err := s.deleteList(id, userID)

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

		var body UpdateListRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := s.updateList(id, userID, body)

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

	r.GET("/lists", func(c *gin.Context) {
		// userID := c.MustGet("userID").(string)
		username := c.Param("username")
		user, err := users.GetUserByUsername(username)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}

		format := c.DefaultQuery("fmt", "preview")
		pageStr := c.DefaultQuery("page", "1")

		val, err := strconv.ParseUint(pageStr, 10, 64)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
			return
		}

		page := uint(val)

		next := cmn.NextPage(c.Request.URL)
		switch format {
		case "preview":
			res, err := s.getListsPreviewByUser(user.ID, &Settings{Page: page})
			res.Next = next

			if err != nil {
				cmn.HandleError(c, err)
				return
			}

			c.IndentedJSON(http.StatusOK, res)
			return

		case "full":
			res, err := s.getListsByUser(user.ID, page)
			res.Next = next

			if err != nil {
				cmn.HandleError(c, err)
				return
			}

			c.IndentedJSON(http.StatusOK, res)
			return
		}

		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid fmt query"})
	})

	r.GET("/lists/recent", func(c *gin.Context) {
		username := c.Param("username")
		user, err := users.GetUserByUsername(username)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}

		res, err := s.getListsPreviewByUser(user.ID, &Settings{Limit: 4})
		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		c.IndentedJSON(http.StatusOK, res)
	})
}
