package property

import "github.com/riazahmedshah/go-booking/internal/server"

type PropertyHandler struct {
	server          *server.Server
	propertyService *PropertyService
}

func NewPropertyHandler(server *server.Server, propertyService *PropertyService) *PropertyHandler {
	return &PropertyHandler{
		server:          server,
		propertyService: propertyService,
	}
}
