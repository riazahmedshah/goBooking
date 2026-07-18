package property

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/go-booking/internal/server"
)

type PropertyHandler struct {
	server          *server.Server
	propertyService *PropertyService
}

func NewPropertyHandler(server *server.Server, propertyService *PropertyService) *PropertyHandler {
	return &PropertyHandler{
		server:          server,
		propertyService: propertyService,
	}
}

func (p *PropertyHandler) CreateProperty(c echo.Context) error {
	var payload CreatePropertyPayload

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	property, err := p.propertyService.CreateProperty(c.Request().Context(), 123, &payload)
	if err != nil {
		slog.Error("failed to create property", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(201, property)
}
