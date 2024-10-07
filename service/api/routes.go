package api

import (
	"example_project/service/api/health"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func (app *CoreAppStruct) routes() http.Handler {
	tokenAuth = jwtauth.New("HS256", []byte(app.Config.SecretKey), nil) // replace with secret key

	mainRouter := chi.NewRouter()
	mainRouter.NotFound(app.NotFound)
	mainRouter.MethodNotAllowed(app.MethodNotAllowed)

	// assigns a unique uuid.v4 in each request
	mainRouter.Use(app.RequestID)

	// logs start and end of a http request
	mainRouter.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  app.Logger,
		NoColor: false,
	}))

	mainRouter.Use(app.RecoverPanic)
	mainRouter.Use(app.Cors)

	healthRoute := health.Routes(&health.HealthAppStruct{Application: app.Application})

	mainRouter.Mount("/api", healthRoute)

	return mainRouter
}
