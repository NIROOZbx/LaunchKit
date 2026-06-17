package apperrors

import "errors"

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrNotFound      = errors.New("not found")
	ErrValidation    = errors.New("validation error")
	ErrInternal      = errors.New("internal server error")
	ErrConflict      = errors.New("conflict")
)
