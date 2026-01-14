package errors

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidInput    = errors.New("invalid input data") // Наша новая метка 🏷️
)
