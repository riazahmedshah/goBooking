package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/go-booking/internal/handler"
	"github.com/riazahmedshah/go-booking/internal/middleware"
)

func registerPropertyRoutes(r *echo.Group, h *handler.Handler, middlewares *middleware.Middlewares) {
	// r.GET("/property", h.PropertyHandler.GetProperty)
	// r.GET("/property/:id", h.PropertyHandler.GetPropertyById)
	property := r.Group("/property")
	property.Use(middlewares.Auth.RequireAuth())
	property.POST("/", h.PropertyHandler.CreateProperty, middlewares.Auth.RequireRole("host"))
}
