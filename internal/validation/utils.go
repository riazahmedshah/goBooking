package validation

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally wrap or format the error here
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
