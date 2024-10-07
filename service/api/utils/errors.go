package utils

import (
	"example_project/internal/response"
	"fmt"
	"net/http"
)

func (app *Application) ErrorMessage(w http.ResponseWriter, r *http.Request, errorCode response.HTTPResponse, errors any) {
	errorCode.JSON(w, r, nil, errors)
}

func (app *Application) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.Logger.WithID(app.GetRequestID(r)).Error(err.Error())
	ServerError.JSON(w, r, nil, nil)
}

func (app *Application) NotFound(w http.ResponseWriter, r *http.Request) {
	ApiNotFound.JSON(w, r, nil, nil)
}

func (app *Application) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	errorMessage := map[string]string{"details": fmt.Sprintf("%v method not allowed", r.Method)}
	MethodNotAllowed.JSON(w, r, nil, errorMessage)
}

func (app *Application) BadRequest(w http.ResponseWriter, r *http.Request, errors any) {
	app.ErrorMessage(w, r, InvalidRequest, errors)
}
