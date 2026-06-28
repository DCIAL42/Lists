package main

import (
	"net/http"
	"time"

	"github.com/DCIAL42/media/internals/client"
	"github.com/DCIAL42/media/internals/music"
	"github.com/DCIAL42/media/internals/search"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	clients := []*client.Client{music.NewMusicClient(httpClient)}

	s := search.NewSearchService(clients...)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/search", s.Search)

	r.Run()
}
