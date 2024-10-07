package utils

import (
	"encoding/json"
	"net/http"
)

type CommonResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Lang    string      `json:"lang"`
	Data    interface{} `json:"data"`
}

type ResponseState struct {
	StatusCode int
	Code       string
	MessageEn  string
	MessageBn  string
}

func (rs ResponseState) CommonResponse(lang string, data interface{}) CommonResponse {
	var message string
	if lang == "bn" {
		message = rs.MessageBn
	} else {
		message = rs.MessageEn
	}
	response := CommonResponse{
		Code:    rs.Code,
		Message: message,
		Lang:    lang,
	}
	if data != nil {
		response.Data = data
	} else {
		response.Data = make(map[string]string)
	}
	return response
}

func (rs ResponseState) WriteToResponse(w http.ResponseWriter, lang string, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rs.StatusCode)
	return json.NewEncoder(w).Encode(rs.CommonResponse(lang, data))
}

var INVALID_REQUEST = ResponseState{
	StatusCode: http.StatusBadRequest,
	Code:       "PAYMENT_TOKEN_INVALID_REQUEST",
	MessageEn:  "token request unsuccessful due to invalid request",
	MessageBn:  "token request unsuccessful due to invalid request",
}

var TOKEN_EXISTS = ResponseState{
	StatusCode: http.StatusBadRequest,
	Code:       "PAYMENT_TOKEN_ALREADY_EXISTS",
	MessageEn:  "token init failed as token already exists",
	MessageBn:  "token init failed as token already exists",
}

var CREATE_TOKEN_SUCCESS = ResponseState{
	StatusCode: http.StatusCreated,
	Code:       "PAYMENT_TOKEN_CREATE_SUCCESS",
	MessageEn:  "token creation successful",
	MessageBn:  "token creation successful",
}

var CREATE_TOKEN_FAILED = ResponseState{
	StatusCode: http.StatusBadRequest,
	Code:       "PAYMENT_TOKEN_CREATE_FAILED",
	MessageEn:  "token creation failed",
	MessageBn:  "token creation failed",
}

var TOKEN_NOT_FOUND = ResponseState{
	StatusCode: http.StatusNotFound,
	Code:       "PAYMENT_TOKEN_ACQUIRE_FAILED",
	MessageEn:  "token not found",
	MessageBn:  "token not found",
}

var TOKEN_REVOKE_FAILED = ResponseState{
	StatusCode: http.StatusBadRequest,
	Code:       "PAYMENT_TOKEN_REVOKE_FAILED",
	MessageEn:  "token revoke unsuccessful",
	MessageBn:  "token revoke unsuccessful",
}

var TOKEN_REVOKE_SUCCESS = ResponseState{
	StatusCode: http.StatusOK,
	Code:       "PAYMENT_TOKEN_REVOKED",
	MessageEn:  "token revoke successful",
	MessageBn:  "token revoke successful",
}

var TOKEN_REFRESH_FAILED = ResponseState{
	StatusCode: http.StatusBadRequest,
	Code:       "PAYMENT_TOKEN_REFRESH_FAILED",
	MessageEn:  "token refresh unsuccessful",
	MessageBn:  "token refresh unsuccessful",
}

var TOKEN_REFRESH_SUCCESS = ResponseState{
	StatusCode: http.StatusOK,
	Code:       "PAYMENT_TOKEN_REFRESH_SUCCESS",
	MessageEn:  "token refresh successful",
	MessageBn:  "token refresh successful",
}

var TOKEN_ACQUIRE_FAILED = ResponseState{
	StatusCode: http.StatusBadRequest,
	Code:       "PAYMENT_TOKEN_ACQUIRE_FAILED",
	MessageEn:  "token acquire unsuccessful",
	MessageBn:  "token acquire unsuccessful",
}

var TOKEN_ALREADY_ACQUIRED = ResponseState{
	StatusCode: http.StatusBadRequest,
	Code:       "PAYMENT_TOKEN_ALREADY_ACQUIRED",
	MessageEn:  "token acquire unsuccessful",
	MessageBn:  "token acquire unsuccessful",
}

var TOKEN_ACQUIRE_SUCCESS = ResponseState{
	StatusCode: http.StatusOK,
	Code:       "PAYMENT_TOKEN_ACQUIRE_SUCCESS",
	MessageEn:  "token acquire successful",
	MessageBn:  "token acquire successful",
}

var INVALID_TOKEN = ResponseState{
	StatusCode: http.StatusUnauthorized,
	Code:       "INVALID_AUTH_TOKEN",
	MessageEn:  "please provide a valid jwt token",
}

var INVALID_PAYMENT_TOKEN = ResponseState{
	StatusCode: http.StatusBadRequest,
	Code:       "INVALID_PAYMENT_TOKEN",
	MessageEn:  "please provide a valid payment token for the customer",
}

var TOKEN_PAYMENT_FAILED = ResponseState{
	StatusCode: http.StatusBadRequest,
	Code:       "TOKEN_PAYMENT_FAILED",
	MessageEn:  "payment failed, please try again",
}

var REFRESH_INACTIVE_TOKEN = ResponseState{
	StatusCode: http.StatusBadRequest,
	Code:       "INACTIVE_TOKEN_REFRESH",
	MessageEn:  "only active tokens can be refreshed",
}

var TRANSACTION_CREATE_FAILED = ResponseState{
	StatusCode: http.StatusBadRequest,
	Code:       "TRANSACTION_CREATE_FAILED",
	MessageEn:  "failed to create a transaction",
	MessageBn:  "",
}
