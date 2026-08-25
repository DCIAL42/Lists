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
	protected.Use(middleware.RequireUser())

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
		count, err := s.getListCount(List{})

		if err != nil {
			cmn.HandleError(c, err)
			return
		}

		res := ListsResponse{
			lists,
			next,
			page,
			count,
		}

		c.IndentedJSON(http.StatusOK, res)
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
		count, err := s.getListCount(List{UserID: user.ID})
		if err != nil {
			cmn.HandleError(c, err)
			return
		}
		switch format {
		case "preview":
			lists, err := s.getListsPreviewByUser(user.ID, page)

			if err != nil {
				cmn.HandleError(c, err)
				return
			}

			c.IndentedJSON(http.StatusOK, ListsPreviewResponse{
				lists,
				next,
				page,
				count,
			})
			return

		case "full":
			lists, err := s.getListsByUser(user.ID, page)

			if err != nil {
				cmn.HandleError(c, err)
				return
			}

			c.IndentedJSON(http.StatusOK, ListsResponse{
				lists,
				next,
				page,
				count,
			})
			return
		}

		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid fmt query"})
	})
}
