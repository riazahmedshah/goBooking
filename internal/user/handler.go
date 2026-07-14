package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/go-booking/internal/server"
)

type UserHandler struct {
	server      *server.Server
	userService *UserService
}

func NewUserHandler(server *server.Server, us *UserService) *UserHandler {
	return &UserHandler{
		server:      server,
		userService: us,
	}
}

func (uh *UserHandler) CreateUser(c echo.Context) error {
	var userPayload CreateUserPayload
	if err := c.Bind(&userPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}
	if err := c.Validate(&userPayload); err != nil {
		return err // Returns the formatted HTTP 400 bad request error
	}
	err := uh.userService.CreateUser(c.Request().Context(), &userPayload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "user created successfully"})
}

func (uh *UserHandler) Login(c echo.Context) error {
	var loginPayload LoginPayload
	if err := c.Bind(&loginPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}
	if err := c.Validate(&loginPayload); err != nil {
		return err // Returns the formatted HTTP 400 bad request error
	}
	result, err := uh.userService.Login(c.Request().Context(), &loginPayload)
	if err != nil {
		if err.Error() == "invalid email or password" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "an unexpected error occurred")
	}
	return c.JSON(http.StatusOK, result)
}
