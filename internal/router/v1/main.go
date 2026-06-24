package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/go-booking/internal/handler"
)

func Registerv1Routes(router *echo.Group, h *handler.Handler) {
	// Register your v1 routes here
	registerUserRoutes(router, h)
}
