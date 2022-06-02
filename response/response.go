package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// global
var gsm *statusMapping

type statusMapping struct {
	mapping map[string]int
}

func newStatusMapping() *statusMapping {
	return &statusMapping{
		mapping: make(map[string]int),
	}
}

func (sm *statusMapping) set(status string, code int) {
	sm.mapping[status] = code
}

// get will return code by set status, it will return 500 (internal server error) when status is not already mapped.
func (sm *statusMapping) get(status string) (code int) {
	code, ok := sm.mapping[status]
	if ok {
		return code
	}
	return http.StatusInternalServerError
}

func init() {

	if gsm == nil {
		gsm := newStatusMapping()
		gsm.set(StatusOK, http.StatusOK)
		gsm.set(StatusCreated, http.StatusCreated)
		gsm.set(StatusNotFound, http.StatusNotFound)
		gsm.set(StatusRequestTimeout, http.StatusRequestTimeout)
		gsm.set(StatusInsufficientBalance, http.StatusForbidden) // register custom status
		// register more if any
	}

}

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

	resp := response{
		err:            err,
		httpStatusCode: gsm.get(status),
		Status:         status,
		Message:        m,
		Data:           data,
		Meta:           meta,
	}

	return resp
}
