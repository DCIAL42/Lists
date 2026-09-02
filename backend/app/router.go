package app

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DCIAL42/lists/client/movies"
	"github.com/DCIAL42/lists/client/music"
	"github.com/DCIAL42/lists/client/search"
	"github.com/DCIAL42/lists/cmn"
	"github.com/DCIAL42/lists/lists"
	"github.com/DCIAL42/lists/media"
	"github.com/DCIAL42/lists/middleware"
	"github.com/DCIAL42/lists/social/follow"
	"github.com/DCIAL42/lists/social/like"
	"github.com/DCIAL42/lists/social/rating"
	"github.com/DCIAL42/lists/tracking"
	"github.com/DCIAL42/lists/users"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupRouter() (*gin.Engine, error) {
	godotenv.Load()

	debug, ok := os.LookupEnv("DEBUG")

	if ok && strings.ToUpper(debug) == "TRUE" {
		slog.Info("Entering debug mode")
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://lists-frontend-lovat.vercel.app",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(location.Default())

	r.Use(middleware.WrapClerkMiddleware(clerkhttp.WithHeaderAuthorization()))

	api := r.Group("/api")

	db, err := cmn.InitDB(
		&lists.List{},
		&lists.ListItem{},
		&cmn.TrackingItem{},
		&follow.Follow{},
		&cmn.Media{},
		&music.Album{},
		&movies.Movie{},
		&cmn.Like{},
		&cmn.Rating{},
	)

	if err != nil {
		return nil, err
	}

	clients := map[cmn.MediaType]cmn.Client{
		cmn.TypeAlbum: music.NewClient(httpClient, db),
		cmn.TypeMovie: movies.NewClient(httpClient, db),
	}

	userGroup := api.Group("/users")
	userService := users.NewService(db)
	userService.SetupRoutes(userGroup)
	userService.SetupUserRoutes(api.Group("/:username"))

	searchGroup := api.Group("/search")
	searchService := search.NewService(clients, db)
	searchService.SetupRoutes(searchGroup)

	likeGroup := api.Group("/like")
	likeService := like.NewService(db)
	likeService.SetupRoutes(likeGroup)

	followGroup := api.Group("/follow")
	followService := follow.NewService(db)
	followService.SetupRoutes(followGroup)
	followService.SetupUserRoutes(api.Group("/:username"))

	listGroup := api.Group("/lists")
	listService := lists.NewService(db, likeService, clients)
	listService.SetupRoutes(listGroup)
	listService.SetupUserRoutes(api.Group("/:username"))

	trackingGroup := api.Group("/tracking")
	trackingService := tracking.NewService(db, clients)
	trackingService.SetupRoutes(trackingGroup)

	ratingGroup := api.Group("/rating")
	ratingService := rating.NewService(db, clients)
	ratingService.SetupRoutes(ratingGroup)
	ratingService.SetupUserRoutes(api.Group("/:username"))

	mediaGroup := api.Group("/media")
	mediaService := media.NewService(db, clients)
	mediaService.SetupRoutes(mediaGroup)
	mediaService.SetupUserRoutes(api.Group("/:username"))

	return r, nil
}
