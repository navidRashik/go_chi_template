package health

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes(router chi.Router, app *HealthAppStruct) http.Handler {
	// router := chi.NewRouter()

	router.Get("/", app.status)
	router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return router
}
