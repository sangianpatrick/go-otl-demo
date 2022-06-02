package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response interface {
	Error() error
	WriteJSON(w http.ResponseWriter) error
	JSONByte() []byte
}

type response struct {
	err            error
	httpStatusCode int
	Status         string      `json:"status"`
	Message        *string     `json:"message,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	Meta           interface{} `json:"meta,omitempty"`
}

func (r response) WriteJSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.httpStatusCode)
	return json.NewEncoder(w).Encode(r)
}

func (r response) Error() error {
	return r.err
}

func (r response) JSONByte() []byte {
	b, _ := json.Marshal(r)
	return b
}

func ResponseSuccess(status string, data, meta interface{}, message string) Response {
	var m *string
	var code int = http.StatusOK

	if message != "" {
		m = &message
	}

	if status == StatusCreated {
		code = http.StatusCreated
	}

	resp := response{
		httpStatusCode: code,
		Status:         status,
		Message:        m,
		Data:           data,
		Meta:           meta,
	}

	return resp
}

func ResponseError(status string, err error, data, meta interface{}, message string) Response {
	var m *string

	if err == nil {
		err = fmt.Errorf("unexpected error")
	}

	if message != "" {
		m = &message
	}

	resp := response{
		err:            err,
		httpStatusCode: getCodeByStatus(status),
		Status:         status,
		Message:        m,
		Data:           data,
		Meta:           meta,
	}

	return resp
}

func getCodeByStatus(status string) int {
	switch status {
	case StatusOK:
		return http.StatusOK
	case StatusCreated:
		return http.StatusCreated
	case StatusNotFound:
		return http.StatusNotFound
	case StatusInternalServerError:
		return http.StatusInternalServerError
	case StatusForbidden:
		return http.StatusForbidden
	case StatusRequestTimeout:
		return http.StatusRequestTimeout
	case StatusBadGateway:
		return http.StatusBadGateway
	case StatusBadRequest:
		return http.StatusBadRequest
	case StatusNotImplemented:
		return http.StatusNotImplemented
	case StatusInsufficientBalance:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
