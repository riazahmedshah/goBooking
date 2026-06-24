package validation

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// This magic snippet makes the validator look at `json:"fieldName"` tags!
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{Validator: v}
}

func (cv *CustomValidator) Validate(i any) error {
	err := cv.Validator.Struct(i)
	if err == nil {
		return nil
	}

	// Check if the error is a collection of validation errors
	if castedErrors, ok := err.(validator.ValidationErrors); ok {
		errorResponse := make(map[string]string)

		for _, fieldErr := range castedErrors {
			// fieldErr.Field() will now give you "firstName" instead of "FirstName"
			fieldName := fieldErr.Field()

			// Customize the message based on the tag that failed
			switch fieldErr.Tag() {
			case "required":
				errorResponse[fieldName] = "This field is required"
			case "email":
				errorResponse[fieldName] = "Invalid email format"
			case "max":
				errorResponse[fieldName] = "Value exceeds maximum allowed length of " + fieldErr.Param()
			case "oneof":
				errorResponse[fieldName] = "Must be one of the following: " + fieldErr.Param()
			default:
				errorResponse[fieldName] = "Invalid value (failed " + fieldErr.Tag() + " restriction)"
			}
		}

		// Return a structured 400 Bad Request
		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "validation failed",
			"errors":  errorResponse,
		})
	}

	// Fallback for any other unexpected errors
	return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
}
