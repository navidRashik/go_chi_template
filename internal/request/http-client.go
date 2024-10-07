package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type CLient struct {
	HttpClient http.Client
	UseMock    bool
	ApiKey     string
	BaseUrl    string
}

func (client *CLient) SendRequest(method string, path string, token string, payload interface{}) (int, []byte, error) {
	// -1, failed to marshal payload into json
	retStatus := -1
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return retStatus, nil, err
	}

	// -2, failed to build a new request
	retStatus--
	requestBody := bytes.NewBuffer(payloadBytes)
	request, err := http.NewRequest(method, path, requestBody)
	if err != nil {
		return retStatus, nil, err
	}

	// -3, A network error is happened
	retStatus--
	request.Header.Add(ContentTypeHeader, ContentTypeValue)
	request.Header.Add(AuthHeader, token)
	// fmt.Printf("%v %v %v\n", request.Method, request.URL, request.Header)

	if client.UseMock {
		request.Header.Add(ApiKey, client.ApiKey)
	}

	// status, return status got from response
	response, err := client.HttpClient.Do(request)
	if err != nil {
		return retStatus, nil, err
	}
	defer response.Body.Close()

	retStatus = response.StatusCode
	body, err := io.ReadAll(response.Body)
	return retStatus, body, err
}
