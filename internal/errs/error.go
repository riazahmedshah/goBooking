package errs

import "fmt"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string, internalErr error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     internalErr,
	}
}
