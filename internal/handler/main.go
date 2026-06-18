package handler

import (
	"github.com/riazahmedshah/go-booking/internal/property"
	"github.com/riazahmedshah/go-booking/internal/server"
	"github.com/riazahmedshah/go-booking/internal/service"
)

type Handler struct {
	PropertyHandler *property.PropertyHandler
}

func NewHandler(server *server.Server, s *service.Service) *Handler {
	propertyHandler := property.NewPropertyHandler(server, s.PropertyService)
	return &Handler{
		PropertyHandler: propertyHandler,
	}
}
