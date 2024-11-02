package errors

import (
	"fmt"
	"github.com/guoxiaopeng875/wallet/internal/pkg/errors/code"
)

var (
	InvalidArgs         = New(code.InvalidArgs, "invalid arguments")
	InsufficientBalance = New(code.InvalidArgs, "insufficient balance")
	RecordNotFound      = New(code.NotFound, "record not found")
	InternalDB          = New(code.InternalServer, "database unknown error")
)

func New(code int, message string) *Error {
	return &Error{
		Status: Status{
			Code:    code,
			Message: message,
		},
	}
}

// Error is a status error.
type Error struct {
	Status
	cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: message = %s  cause = %v", e.Message, e.cause)
}

// WithCause with the underlying cause of the error.
func (e *Error) WithCause(cause error) *Error {
	err := Clone(e)
	err.cause = cause
	return err
}

// Clone deep clone error to a new error.
func Clone(err *Error) *Error {
	return &Error{
		cause: err.cause,
		Status: Status{
			Code:    err.Code,
			Message: err.Message,
		},
	}
}
