package apperrors

import "errors"

var (
	ErrUnauthorized  = errors.New("unauthorized: invalid or missing session")
	ErrForbidden     = errors.New("forbidden")
	ErrNotFound      = errors.New("not found")
	ErrValidation    = errors.New("validation error")
	ErrInternal      = errors.New("internal server error")
	ErrConflict      = errors.New("conflict")
	ErrRefreshLockHeld = errors.New("refresh lock held: another request is already processing a refresh")
)
