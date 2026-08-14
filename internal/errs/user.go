package errs

import "net/http"

var (
	ErrDuplicateEmail = New(
		http.StatusConflict,
		"email already alreadt taken, please use different email",
		nil,
	)

	ErrUserNotFound = New(
		http.StatusNotFound,
		"user not found",
		nil,
	)

	ErrInvalidPassword = New(
		http.StatusUnauthorized,
		"invalid password",
		nil,
	)
)
