package apperr

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
