package repository

import (
	"github.com/riazahmedshah/go-booking/internal/property"
	"github.com/riazahmedshah/go-booking/internal/server"
	"github.com/riazahmedshah/go-booking/internal/user"
)

type Repositories struct {
	UserRepository     *user.UserRepository
	PropertyRepository *property.PropertyRepository
}

func NewRepositories(s *server.Server) *Repositories {
	userRepo := user.NewUserRepository(s)
	propertyRepo := property.NewPropertyRepository(s)
	return &Repositories{
		UserRepository:     userRepo,
		PropertyRepository: propertyRepo,
	}
}
