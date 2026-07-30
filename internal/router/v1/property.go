package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/go-booking/internal/handler"
)

func registerPropertyRoutes(r *echo.Group, h *handler.Handler) {
	property := r.Group("/property")

	property.POST("", h.PropertyHandler.CreateProperty)
}
