package handler

import (
	"github.com/riazahmedshah/go-booking/internal/property"
	"github.com/riazahmedshah/go-booking/internal/server"
	"github.com/riazahmedshah/go-booking/internal/service"
	"github.com/riazahmedshah/go-booking/internal/user"
)

type Handler struct {
	UserHandler     *user.UserHandler
	PropertyHandler *property.PropertyHandler
}

func NewHandler(server *server.Server, service *service.Service) *Handler {
	userHandler := user.NewUserHandler(server, service.UserService)
	propertyHandler := property.NewPropertyHandler(server, service.PropertyService)
	return &Handler{
		UserHandler:     userHandler,
		PropertyHandler: propertyHandler,
	}
}
