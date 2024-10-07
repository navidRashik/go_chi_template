package response

import (
	"encoding/json"
	"net/http"
)

type HTTPResponse struct {
	Status    int
	Code      string
	MessageEn string
	MessageBn string
}

type responseStruct struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Lang    string      `json:"lang"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors,omitempty"`
}

func (rs HTTPResponse) prepareResponse(lang string, data interface{}, errors interface{}) responseStruct {
	var message string
	if lang == "bn" {
		message = rs.MessageBn
	} else {
		message = rs.MessageEn
	}
	response := responseStruct{
		Code:    rs.Code,
		Message: message,
		Lang:    lang,
	}
	if data != nil {
		response.Data = data
	} else {
		response.Data = make(map[string]string)
	}

	if errors != nil {
		response.Errors = errors
	} else {
		response.Errors = nil
	}
	return response
}

func (rs *HTTPResponse) JSON(w http.ResponseWriter, r *http.Request, data any, errors any) error {
	return rs.JSONWithHeaders(w, r, nil, data, errors)
}

func (rs *HTTPResponse) JSONWithHeaders(w http.ResponseWriter, r *http.Request, headers http.Header, data any, errors any) error {
	response := rs.prepareResponse(getLanguage(r), data, errors)
	responseJson, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-Id", r.Header.Get("X-Request-Id"))

	responseJson = append(responseJson, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.WriteHeader(rs.Status)
	w.Write(responseJson)
	return nil
}
