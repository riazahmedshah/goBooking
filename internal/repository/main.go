package repository

import (
	"github.com/riazahmedshah/go-booking/internal/property"
	"github.com/riazahmedshah/go-booking/internal/server"
)

type Repositories struct {
	PropertyRepo *property.PropertyRepository
}

func NewRepositories(s *server.Server) *Repositories {
	propertyRepo := property.NewPropertyRepository(s)
	return &Repositories{
		PropertyRepo: propertyRepo,
	}
}
