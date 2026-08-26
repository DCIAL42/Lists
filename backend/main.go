package main

import (
	"log"
	"os"

	"github.com/DCIAL42/lists/app"
)

func main() {
	r, err := app.SetupRouter()

	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "38247"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
