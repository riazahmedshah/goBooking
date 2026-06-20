package booking

import "github.com/go-playground/validator/v10"

type CreateBookingPayload struct {
	UserID     *int     `json:"userId" validate:"required"`
	PropertyID *int     `json:"propertyId" validate:"required"`
	TotalPrice *float64 `json:"totalPrice" validate:"required,gt=0"`
	Status     *string  `json:"status" validate:"required,oneof=pending confirmed cancelled"`
}

func (p *CreateBookingPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type ConfirmBookingPayload struct {
	Status *string `json:"status" validate:"required,oneof=pending confirmed cancelled"`
}

func (p *ConfirmBookingPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type CreateIdempotencyKeyPayload struct {
	Key *string `json:"key" validate:"required"`
}

func (p *CreateIdempotencyKeyPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type FinalizeIdempotencyKeyPayload struct {
	IsFinalized *bool `json:"isFinalized" validate:"required"`
}
