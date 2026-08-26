package main

import (
	"github.com/DCIAL42/lists/app"
)

func main() {
	r := app.SetupRouter()

	r.Run(":8080")
}
