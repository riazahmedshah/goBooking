package service

import (
	"github.com/riazahmedshah/go-booking/internal/property"
	"github.com/riazahmedshah/go-booking/internal/repository"
	"github.com/riazahmedshah/go-booking/internal/server"
)

type Service struct {
	PropertyService *property.PropertyService
}

func NewService(server *server.Server, repository *repository.Repositories) (*Service, error) {
	propertyService := property.NewPropertyService(server, repository.PropertyRepo)
	return &Service{
		PropertyService: propertyService,
	}, nil
}
