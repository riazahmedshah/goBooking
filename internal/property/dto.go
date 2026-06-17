package property

import "github.com/go-playground/validator/v10"

type CreatePropertyPayload struct {
	Title     string  `json:"title" validate:"required,min=1,max=255"`
	SubTitle  *string `json:"subTitle" validate:"omitempty,max=1000"`
	Image     *string `json:"image" validate:"omitempty"`
	AddressID int     `json:"addressId" validate:"required"`
	MaxGuests *int    `json:"maxGuest" validate:"omitempty,min=1"`
}

func (p *CreatePropertyPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// TODO: other payloads...
