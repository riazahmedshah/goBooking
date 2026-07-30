package property

import (
	"context"
	"log/slog"

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

func (ps *PropertyService) CreateProperty(ctx context.Context, hostID string, payload *CreatePropertyPayload) (*Property, error) {
	property, err := ps.propertyRepo.Createproperty(ctx, hostID, payload)
	if err != nil {
		slog.Error("database failure during property creation", "error", err)
		return nil, err
	}

	return property, nil
}

func (ps *PropertyService) GetAllProperties(ctx context.Context) ([]*Property, error) {
	properties, err := ps.propertyRepo.GetAllProperties(ctx)
	if err != nil {
		slog.Error("database failure during property retrieval", "error", err)
		return nil, err
	}

	return properties, nil
}

func (ps *PropertyService) GetPropertyByID(ctx context.Context, propertyID string) (*Property, error) {
	property, err := ps.propertyRepo.GetPropertyByID(ctx, propertyID)
	if err != nil {
		slog.Error("database failure during property retrieval by ID", "error", err)
		return nil, err
	}

	return property, nil
}
