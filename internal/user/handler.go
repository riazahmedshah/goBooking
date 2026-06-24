package user

import (
	"fmt"
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

func (*UserHandler) CreateUser(c echo.Context) error {
	var userPayload CreateUserPayload
	if err := c.Bind(&userPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}
	if err := c.Validate(&userPayload); err != nil {
		return err // Returns the formatted HTTP 400 bad request error
	}
	fmt.Println("All validation passed!")
	return c.NoContent(http.StatusCreated)
}
