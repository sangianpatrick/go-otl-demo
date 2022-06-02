package exception

import "fmt"

var (
	ErrNotFound            = fmt.Errorf("error: desirable data is not found")
	ErrInternalServer      = fmt.Errorf("error: internal server error")
	ErrBadGateway          = fmt.Errorf("error: bad gateway")
	ErrInsufficientBalance = fmt.Errorf("error: insufficient balance")
	ErrNotImplemented      = fmt.Errorf("error: not implemented")
	ErrBadRequest          = fmt.Errorf("error: bad request")
)
