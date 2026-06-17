package property

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/riazahmedshah/go-booking/internal/server"
)

type PropertyRepository struct {
	server *server.Server
}

func NewPropertryRepository(server *server.Server) *PropertyRepository {
	return &PropertyRepository{server: server}
}

func (p *PropertyRepository) Createproperty(ctx context.Context, hostID int, payload *CreatePropertyPayload) (*Property, error) {
	stmt := `
		INSERT INTO properties(
			host_id,
			title,
			sub_title,
			image,
			address_id
			max_guest
		)
		VALUES (
			@host_id,
			@title,
			@sub_title,
			@image,
			@address_id,
			@max_guest
		)
		RETURNING
	`

	rows, err := p.server.DB.Query(ctx, stmt, pgx.NamedArgs{
		"host_id":    hostID,
		"title":      payload.Title,
		"sub_title":  payload.SubTitle,
		"image":      payload.Image,
		"address_id": payload.AddressID,
		"max_guest":  payload.MaxGuests,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute create todo query for host_id=%v title=%s: %w", hostID, payload.Title, err)
	}

	propertyItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Property])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:todos for host_id=%v title=%s: %w", hostID, payload.Title, err)
	}

	return &propertyItem, nil
}
