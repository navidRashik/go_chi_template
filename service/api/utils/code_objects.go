package utils

import (
	"net/http"

	"example_project/internal/response"
)

var HealthResponse = response.HTTPResponse{
	Status:    http.StatusOK,
	Code:      "LIVE_OK_200",
	MessageEn: "alive and well",
	MessageBn: "",
}

var ApiNotFound = response.HTTPResponse{
	Status:    http.StatusNotFound,
	Code:      "EPNF_404",
	MessageEn: "The requested resource could not be found",
	MessageBn: "",
}

var MethodNotAllowed = response.HTTPResponse{
	Status:    http.StatusMethodNotAllowed,
	Code:      "EPMNA_405",
	MessageEn: "Method not supported for this resource",
	MessageBn: "",
}

var RequestAccepted = response.HTTPResponse{
	Status:    http.StatusAccepted,
	Code:      "EPRR_202",
	MessageEn: "Request received successfully",
	MessageBn: "Request received successfully",
}

var InvalidRequest = response.HTTPResponse{
	Status:    http.StatusBadRequest,
	Code:      "EPIR_400",
	MessageEn: "EPta valiEPtion failed",
	MessageBn: "EPta valiEPtion failed",
}

var ServerError = response.HTTPResponse{
	Status:    http.StatusInternalServerError,
	Code:      "EPISE_500",
	MessageEn: "An unexpected error happened in the server",
	MessageBn: "An unexpected error happened in the server",
}

var SuccessResponse = response.HTTPResponse{
	Status:    http.StatusOK,
	Code:      "EPS_200",
	MessageEn: "Success",
	MessageBn: "Success",
}

var IssoAuthRequired = response.HTTPResponse{
	Status:    http.StatusUnauthorized,
	Code:      "EPU_401",
	MessageEn: "You are unauthorized to perform this request",
	MessageBn: "You are unauthorized to perform this request",
}

var InvalidCredentials = response.HTTPResponse{
	Status:    http.StatusUnauthorized,
	Code:      "EPIC_400",
	MessageEn: "Login failed due to invalid credentials",
	MessageBn: "Login failed due to invalid credentials",
}

var UserCreateSuccess = response.HTTPResponse{
	Status:    http.StatusOK,
	Code:      "EPCUS_200",
	MessageEn: "User created successfully",
	MessageBn: "User created successfully",
}
