package apperr

import "net/http"

type AppError struct {
	StatusCode int
	Message    string
	Log        string
}

// Error returns error message.
// AppErr satisfies error interface.
func (e AppError) Error() string {
	return e.Message
}

var (
	ServerError      = AppError{http.StatusInternalServerError, "internal server error", "internal server error"}
	ShortURLNotFound = AppError{http.StatusNotFound, "short url not found", ""}
)
