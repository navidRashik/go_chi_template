package health

import (
	"example_project/internal/version"

	utils "example_project/service/api/utils"
	"net/http"
	"time"
)

func (app *HealthAppStruct) status(w http.ResponseWriter, r *http.Request) {
	logger := app.Logger.WithID(app.GetRequestID(r))
	logger.Print("health request received")
	data := map[string]string{
		"status":  "OK",
		"version": version.GetVersion(),
		"now":     time.Now().Format("02-01-2006 15:04"),
	}
	err := utils.HealthResponse.JSON(w, r, data, nil)
	if err != nil {
		logger.Error("error happened during sending response, details: %s", err.Error())
	}
}
