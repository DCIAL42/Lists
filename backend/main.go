package main

import (
	"log"
	"os"

	"github.com/DCIAL42/lists/app"
)

func main() {
	log.Printf("PORT=%q", os.Getenv("PORT"))
	log.Printf("VERCEL=%q", os.Getenv("VERCEL"))
	log.Printf("VERCEL_ENV=%q", os.Getenv("VERCEL_ENV"))

	r, err := app.SetupRouter()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is missing")
	}

	log.Printf("Starting server on :%s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
