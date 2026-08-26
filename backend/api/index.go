package api

import (
	"log"
	"net/http"

	"github.com/DCIAL42/lists/app"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("=== HANDLER ENTERED ===")

	router, err := app.SetupRouter()
	if err != nil {
		log.Printf("=== SETUP ROUTER ERROR: %v ===", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("=== ROUTER CREATED ===")

	router.ServeHTTP(w, r)

	log.Println("=== REQUEST SERVED ===")
}
