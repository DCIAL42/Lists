package api

import (
	"net/http"

	"github.com/DCIAL42/lists/app"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router, err := app.SetupRouter()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	router.ServeHTTP(w, r)
}
