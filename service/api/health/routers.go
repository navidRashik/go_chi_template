package health

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Routes(app *HealthAppStruct) http.Handler {
	mux := chi.NewMux()
	// healthApp := healthStruct{app}
	mux.Get("/health", app.status)
	return mux
}
