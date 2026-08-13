package booking

import (
	"github.com/go-playground/validator/v10"
)

type CreateBookingPayload struct {
	PropertyID *string  `json:"propertyId" validate:"required"`
	TotalPrice *float64 `json:"totalPrice" validate:"required,gt=0"`
	CheckIn    *string  `json:"checkIn" validate:"required,datetime=2006-01-02"`
	CheckOut   *string  `json:"checkOut" validate:"required,datetime=2006-01-02,gtfield=CheckIn"`
}

func (p *CreateBookingPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type ConfirmBookingPayload struct {
	BookingID string `json:"bookingId" validate:"required"`
	Status    string `json:"status" validate:"omitempty,oneof=pending confirmed cancelled"`
}

func (p *ConfirmBookingPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type FinalizeIdempotencyKeyPayload struct {
	IsFinalized *bool `json:"isFinalized" validate:"required"`
}
