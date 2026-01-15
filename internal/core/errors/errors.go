package errors

import (
	"errors"
	"fmt"
)

//type ClientError interface {
//	Error() string
//	IsClientError() bool
//}

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrInvalidInput        = errors.New("invalid input data")
	ErrMethodNotAllowed    = errors.New("method not allowed")
	ErrInternalServerError = errors.New("internal server error")
)

type ValidationError struct {
	Field string
	Value string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid value %s in field %s", e.Value, e.Field)
}
