package property

import "github.com/riazahmedshah/go-booking/internal/server"

type PropertyService struct {
	server       *server.Server
	propertyRepo *PropertyRepository
}

func NewPropertyService(server *server.Server, propertyRepo *PropertyRepository) *PropertyService {
	return &PropertyService{
		server:       server,
		propertyRepo: propertyRepo,
	}
}
