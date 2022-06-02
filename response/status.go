package response

const (
	StatusOK                  string = "OK"
	StatusCreated             string = "CREATED"
	StatusNotFound            string = "NOT_FOUND"
	StatusInternalServerError string = "INTERNAL_SERVER_ERROR"
	StatusForbidden           string = "FORBIDDEN"
	StatusRequestTimeout      string = "REQUEST_TIMEOUT"
	StatusBadGateway          string = "BAD_GATEWAY"
	StatusInsufficientBalance string = "INSUFFICIENT_BALANCE" // custom status
	// add more if any
)
