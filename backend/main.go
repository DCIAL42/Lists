package main

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DCIAL42/media/cmn"
	"github.com/DCIAL42/media/internals/client"
	"github.com/DCIAL42/media/internals/movies"
	"github.com/DCIAL42/media/internals/music"
	"github.com/DCIAL42/media/internals/search"
	"github.com/DCIAL42/media/lists"
	"github.com/DCIAL42/media/tracking"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	debug, ok := os.LookupEnv("DEBUG")
	if ok && strings.ToUpper(debug) == "TRUE" {
		slog.Info("Entering debug mode")
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	musicClient := music.NewMusicClient(httpClient)
	movieClient := movies.NewMovieClient(httpClient)
	clients := map[cmn.MediaType]client.Client{
		cmn.TypeAlbum: musicClient,
		cmn.TypeMovie: movieClient,
	}

	s := search.NewSearchService(clients)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(location.Default())

	// protected := r.Group("/")
	// protected.Use(middleware.WrapClerkMiddleware(clerkhttp.WithHeaderAuthorization()), middleware.RequireUser())

	api := r.Group("/api")

	api.GET("/search", s.Search)

	db, err := cmn.InitDB(&lists.List{}, &lists.ListItem{}, &tracking.TrackingItem{})

	if err != nil {
		panic(err)
	}

	listGroup := api.Group("/lists")
	listService := lists.NewService(db, clients)
	listService.SetupRoutes(listGroup)

	trackingGroup := api.Group("/tracking")
	trackingService := tracking.NewService(db, clients)
	trackingService.SetupRoutes(trackingGroup)

	r.Run(":8080")
}
