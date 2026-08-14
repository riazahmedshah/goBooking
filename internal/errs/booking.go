package errs

import "net/http"

var (
	ErrPropertyHeld = New(
		http.StatusConflict,
		"property is currently held by another booking request, please try again later",
		nil,
	)

	ErrBookingNotFound = New(
		http.StatusNotFound,
		"booking record not found",
		nil,
	)

	ErrBookingInProgress = New(
		http.StatusOK,
		"your booking is already in progress, please check your payments or wait a moment",
		nil,
	)

	ErrDuplicateBooking = New(
		http.StatusConflict,
		"booking is already finalized",
		nil,
	)
)
