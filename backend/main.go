package main

import (
	"log"

	"github.com/DCIAL42/lists/app"
)

func main() {
	r, err := app.SetupRouter()

	if err != nil {
		log.Fatal(err)
	}

	r.Run(":8080")
}
