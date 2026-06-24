package handler

import (
	"github.com/riazahmedshah/go-booking/internal/server"
	"github.com/riazahmedshah/go-booking/internal/service"
	"github.com/riazahmedshah/go-booking/internal/user"
)

type Handler struct {
	UserHandler *user.UserHandler
}

func NewHandler(server *server.Server, service *service.Service) *Handler {
	userHandler := user.NewUserHandler(server, service.UserService)
	return &Handler{
		UserHandler: userHandler,
	}
}
