package service

import (
	"github.com/riazahmedshah/go-booking/internal/repository"
	"github.com/riazahmedshah/go-booking/internal/server"
	"github.com/riazahmedshah/go-booking/internal/user"
)

type Service struct {
	UserService *user.UserService
}

func NewService(server *server.Server, repository *repository.Repositories) (*Service, error) {
	userService := user.NewUserService(server, repository.UserRepository)
	return &Service{
		UserService: userService,
	}, nil
}
