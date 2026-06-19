package property

import (
	"context"

	"github.com/riazahmedshah/go-booking/internal/server"
)

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

func (p *PropertyService) CreateProperty(ctx context.Context, hostID int, payload *CreatePropertyPayload) (*Property, error) {
	property, err := p.propertyRepo.Createproperty(ctx, hostID, payload)
	if err != nil {
		return nil, err
	}

	return property, nil
}
