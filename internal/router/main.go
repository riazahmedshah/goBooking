package router

import (
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/go-booking/internal/handler"
	v1 "github.com/riazahmedshah/go-booking/internal/router/v1"
	"github.com/riazahmedshah/go-booking/internal/server"
	"github.com/riazahmedshah/go-booking/internal/validation"
)

func NewRouter(s *server.Server, h *handler.Handler) *echo.Echo {
	router := echo.New()

	router.Validator = validation.NewCustomValidator()

	// Register your routes here
	v1Group := router.Group("/api/v1")
	v1.Registerv1Routes(v1Group, h)

	return router
}
