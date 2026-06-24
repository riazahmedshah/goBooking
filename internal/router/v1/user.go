package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/go-booking/internal/handler"
)

func registerUserRoutes(r *echo.Group, h *handler.Handler) {
	auth := r.Group("/auth")

	auth.POST("/register", h.UserHandler.CreateUser)
	// auth.POST("/login", )
}
