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

func NewPropertyRepository(server *server.Server) *PropertyRepository {
	return &PropertyRepository{server: server}
}

func (p *PropertyRepository) Createproperty(ctx context.Context, hostID int, payload *CreatePropertyPayload) (*Property, error) {
	stmt := `
		INSERT INTO properties(
			host_id,
			title,
			sub_title,
			image,
			address_id,
			max_guests
		)
		VALUES (
			@host_id,
			@title,
			@sub_title,
			@image,
			@address_id,
			@max_guests
		)
		RETURNING *
	`

	rows, err := p.server.DB.Query(ctx, stmt, pgx.NamedArgs{
		"host_id":    hostID,
		"title":      payload.Title,
		"sub_title":  payload.SubTitle,
		"image":      payload.Image,
		"address_id": payload.AddressID,
		"max_guests": payload.MaxGuests,
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

func (p *PropertyRepository) GetPropertyByID(ctx context.Context, propertyID int) (*Property, error) {
	stmt := `
		SELECT
			id,
			title,
			subtitle,
			image,
			address_id,
			host_id,
			created-at,
			updated_at
		FROM
			properties
		WHERE
			id = @id
		`

	rows, err := p.server.DB.Query(ctx, stmt, pgx.NamedArgs{
		"id": propertyID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute get property by id for property_id %v: %w", propertyID, err)
	}

	propertyItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Property])

	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:properties for property_id=%d: %w", propertyID, err)
	}

	return &propertyItem, nil

}
