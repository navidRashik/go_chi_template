package health

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Routes(app *HealthAppStruct) http.Handler {
	router := chi.NewRouter()

	router.Get("/", app.status)
	router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return router
}
