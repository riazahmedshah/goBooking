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

func (ph *PropertyHandler) CreateProperty(c echo.Context) error {
	userID, _ := c.Get("userID").(string)
	var payload CreatePropertyPayload

	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	property, err := ph.propertyService.CreateProperty(c.Request().Context(), userID, &payload)
	if err != nil {
		slog.Error("failed to create property", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(201, property)
}

func (ph *PropertyHandler) GetAllProperties(c echo.Context) error {
	properties, err := ph.propertyService.GetAllProperties(c.Request().Context())
	if err != nil {
		slog.Error("failed to retrieve properties", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(200, properties)
}

func (ph *PropertyHandler) GetPropertyById(c echo.Context) error {
	propertyID := c.Param("id")
	property, err := ph.propertyService.GetPropertyByID(c.Request().Context(), propertyID)
	if err != nil {
		slog.Error("failed to retrieve property by ID", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	if property == nil {
		return echo.NewHTTPError(http.StatusNotFound, "property not found")
	}
	return c.JSON(200, property)
}

func (ph *PropertyHandler) GetPropertyAvailability(c echo.Context) error {
	propertyID := c.Param("id")
	availability, err := ph.propertyService.GetPropertyAvailability(c.Request().Context(), propertyID)
	if err != nil {
		slog.Error("failed to retrieve property availability", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(200, availability)
}
