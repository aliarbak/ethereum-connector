package errors

import (
	"fmt"
)

type Error struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func New(message string, a ...any) error {
	return &Error{
		StatusCode: 500,
		Message:    fmt.Sprintf(message, a...),
	}
}

func From(err error, message string, a ...any) error {
	return &Error{
		StatusCode: 500,
		Message:    fmt.Sprintf(fmt.Sprintf("%s, err: %s", message, err.Error()), a...),
	}
}

func NewWithStatusCode(statusCode int, message string, a ...any) error {
	return &Error{
		StatusCode: statusCode,
		Message:    fmt.Sprintf(message, a...),
	}
}

func (r *Error) Error() string {
	return fmt.Sprintf("(statusCode=%d) %s", r.StatusCode, r.Message)
}

func FromInvalidInput(err error) error {
	return &Error{
		StatusCode: 400,
		Message:    err.Error(),
	}
}

func InvalidInput(message string, a ...any) error {
	return NewWithStatusCode(400, message, a...)
}

func AlreadyExists(message string, a ...any) error {
	return NewWithStatusCode(409, message, a...)
}

func NotFound(message string, a ...any) error {
	return NewWithStatusCode(404, message, a...)
}
