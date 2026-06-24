package user

import "github.com/go-playground/validator/v10"

type CreateUserPayload struct {
	FirstName string  `json:"firstName" validate:"required,max=255"`
	LastName  *string `json:"lastName" validate:"omitempty,max=255"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6,max=20"`
	Role      *string `json:"role" validate:"omitempty,oneof=user host"`
}

func (payload *CreateUserPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(payload)
}

type CreateHostPayload struct {
	UserID    string `json:"userId" validate:"required,uuid"`
	StateName string `json:"stateName" validate:"required,min=3,max=50"`
	City      string `json:"city" validate:"required,min=1,max=100"`
	Area      string `json:"area" validate:"required,min=1,max=500"`
}
