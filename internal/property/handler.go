package property

import (
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
		return err
	}

	property, err := p.propertyService.CreateProperty(c.Request().Context(), 123, &payload)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "internal server error"})
	}
	return c.JSON(201, property)
}
