package errs

import "net/http"

var (
	ErrPropertyTitleExists = New(
		http.StatusConflict,
		"use different title for property, this title already exists",
		nil,
	)

	ErrPropertyNotFound = New(
		http.StatusNotFound,
		"property not found",
		nil,
	)

	ErrBadUpdateRequest = New(
		http.StatusBadRequest,
		"invalid update fields",
		nil,
	)
)
