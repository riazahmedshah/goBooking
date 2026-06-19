package errs

import "net/http"

func NewBadRequestError(message string, code *string, errors []FieldError, action *Action) *HTTPError {
	formattedCode := MakeUpperCaseWithUnderscores(http.StatusText(http.StatusBadRequest))

	if code != nil {
		formattedCode = *code
	}

	return &HTTPError{
		Code:    formattedCode,
		Message: message,
		Status:  http.StatusBadRequest,
		Errors:  errors,
		Action:  action,
	}
}

func NewNotFoundError(message string, code *string) *HTTPError {
	formattedCode := MakeUpperCaseWithUnderscores(http.StatusText(http.StatusNotFound))

	if code != nil {
		formattedCode = *code
	}

	return &HTTPError{
		Code:    formattedCode,
		Message: message,
		Status:  http.StatusNotFound,
	}
}
