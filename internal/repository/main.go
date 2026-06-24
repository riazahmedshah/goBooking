package repository

import (
	"github.com/riazahmedshah/go-booking/internal/server"
	"github.com/riazahmedshah/go-booking/internal/user"
)

type Repositories struct {
	UserRepository *user.UserRepository
}

func NewRepositories(s *server.Server) *Repositories {
	userRepo := user.NewUserRepository(s)
	return &Repositories{
		UserRepository: userRepo,
	}
}
