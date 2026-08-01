package main

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DCIAL42/media/internals/client"
	"github.com/DCIAL42/media/internals/movies"
	"github.com/DCIAL42/media/internals/music"
	"github.com/DCIAL42/media/internals/search"
	"github.com/DCIAL42/media/lists"
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
	clients := []client.Client{musicClient, movieClient}

	s := search.NewSearchService(clients...)

	r := gin.Default()

	api := r.Group("/api")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(location.Default())

	r.GET("/search", s.Search)

	listGroup := api.Group("/lists")

	listService := lists.NewService(musicClient, movieClient)

	listService.SetupRoutes(listGroup)

	r.Run(":8080")
}
