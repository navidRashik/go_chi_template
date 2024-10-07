package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const AuthHeader = "Authorization"
const ContentTypeHeader = "Content-Type"
const ContentTypeValue = "application/json"
const ApiKey = "x-api-key"
const ApiKeyValue = "PMAK-631d7a70b5dcb2b3087377d4845dc85d736f763"

// PostJson -->
//
// -1 -> json marshaling error,
//
// -2 ->  couldn't create a request error
//
// -3 -> a server side error,
func PostJson(url string, token string, payload interface{}, TimeOut time.Duration) (int, []byte, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), TimeOut)
	defer cancel()
	retStatus := -1
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return retStatus, nil, err
	}
	retStatus--
	requestBody := bytes.NewBuffer(payloadBytes)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, requestBody)
	if err != nil {
		return retStatus, nil, err
	}
	request.Header.Add(ContentTypeHeader, ContentTypeValue)
	request.Header.Add(AuthHeader, token)
	// TODO: Adding api-key for mock response
	request.Header.Add(ApiKey, ApiKeyValue)
	retStatus--
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return retStatus, nil, err
	}
	defer response.Body.Close()
	retStatus = response.StatusCode
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return retStatus, nil, err
	}
	return retStatus, body, nil
}

func IsNetworkError(status int) bool {
	return status == -3
}

func IsBuildRequestError(status int) bool {
	return status == -1 || status == -2
}

func IsSuccess(status int) bool {
	return 200 <= status && status <= 299
}

func IsBadRequestError(status int) bool {
	return 400 <= status && status <= 499
}

func IsRedirect(status int) bool {
	return 300 <= status && status <= 399
}

func IsServerError(status int) bool {
	return 500 <= status && status <= 599
}

func IsInformational(status int) bool {
	return 100 <= status && status <= 199
}

type ErrType int

const (
	ServerError ErrType = iota
	BuildRequestError
	RequestSuccess
	BadRequestError
	RedirectError
	InformationalError
	UnknownError
)

func GetErrorType(status int) ErrType {
	if IsBuildRequestError(status) {
		return BuildRequestError
	} else if IsInformational(status) {
		return InformationalError
	} else if IsRedirect(status) {
		return RedirectError
	} else if IsBadRequestError(status) {
		return BadRequestError
	} else if IsServerError(status) || IsNetworkError(status) {
		return ServerError
	} else if IsSuccess(status) {
		return RequestSuccess
	}
	return UnknownError
}
