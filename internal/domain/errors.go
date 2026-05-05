package domain

import "errors"

var (
	ErrNotFound              = errors.New("resource not found")
	ErrInvalidInput          = errors.New("invalid input")
	ErrConflict              = errors.New("conflict")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrInternalServer        = errors.New("internal server error")
	ErrReservationOverlap    = errors.New("reservation overlaps with existing reservation")
	ErrSportHasCourts        = errors.New("sport has associated courts")
	ErrInvalidDateRange      = errors.New("invalid date range")
	ErrCourtNotFound         = errors.New("court not found")
	ErrSportNotFound         = errors.New("sport not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrReservationNotFound   = errors.New("reservation not found")
	ErrDuplicateSportName    = errors.New("sport name already exists")
	ErrDuplicateEmail        = errors.New("email already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
)

type ValidationError struct {
	Field   string
	Message string
}

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

func NewAppError(code, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}
