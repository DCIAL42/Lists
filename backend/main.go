package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/DCIAL42/media/internals/client"
	"github.com/DCIAL42/media/internals/movies"
	"github.com/DCIAL42/media/internals/music"
	"github.com/DCIAL42/media/internals/search"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	clients := []*client.Client{music.NewMusicClient(httpClient), movies.NewMovieClient(httpClient)}

	s := search.NewSearchService(clients...)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(location.Default())

	r.GET("/search", s.Search)

	r.Run(":8080")
}
