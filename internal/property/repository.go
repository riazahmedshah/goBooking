package property

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/riazahmedshah/go-booking/internal/errs"
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
		return nil, fmt.Errorf("failed to execute create todo query for host_id=%d title=%s: %w", hostID, payload.Title, err)
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

func (p *PropertyRepository) UpdateProperty(ctx context.Context, propertyID int, payload *UpdatePropertyPayload) (*Property, error) {
	stmt := "UPDATE properties SET "

	args := pgx.NamedArgs{
		"id": propertyID,
	}

	setClauses := []string{}

	if payload.SubTitle != nil {
		setClauses = append(setClauses, "sub_title = @sub_title")
		args["sub_title"] = *payload.SubTitle
	}

	if payload.Image != nil {
		setClauses = append(setClauses, "image = @image")
		args["image"] = *payload.Image
	}

	if payload.AddressID != nil {
		setClauses = append(setClauses, "address_id = @address_id")
		args["address_id"] = *payload.AddressID
	}

	if payload.MaxGuests != nil {
		setClauses = append(setClauses, "max_guests = @max_guests")
		args["max_guests"] = *payload.MaxGuests
	}

	if len(setClauses) == 0 {
		return nil, errs.NewBadRequestError("no fields to update", nil, nil, nil)
	}

	stmt += strings.Join(setClauses, ", ")
	stmt += " WHERE id = @id RETURNING *"

	rows, err := p.server.DB.Query(ctx, stmt, args)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	updatedProperty, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Property])

	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:properties for property_id=%d: %w", propertyID, err)
	}

	return &updatedProperty, nil

}

func (p *PropertyRepository) DeleteProperty(ctx context.Context, propertyID int) error {
	stmt := `
		DELETE FROM properties
		WHERE id = @id
	`

	result, err := p.server.DB.Exec(ctx, stmt, pgx.NamedArgs{
		"id": propertyID,
	})

	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if result.RowsAffected() == 0 {
		code := "PROPERTY_NOT_FOUND"
		return errs.NewNotFoundError("property not found", &code)
	}

	return nil
}
