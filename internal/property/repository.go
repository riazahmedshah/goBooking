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

func (pr *PropertyRepository) Createproperty(ctx context.Context, hostID string, payload *CreatePropertyPayload) (*Property, error) {
	stmt := `
		INSERT INTO properties(
			host_id, title, sub_title, max_guests, price
		)
		VALUES (
			@host_id, @title, @sub_title, @max_guests, @price
		)
		RETURNING *
	`

	rows, err := pr.server.DB.Query(ctx, stmt, pgx.NamedArgs{
		"host_id":    hostID,
		"title":      payload.Title,
		"sub_title":  payload.SubTitle,
		"max_guests": payload.MaxGuests,
		"price":      payload.Price,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute create property query for host_id=%s title=%s: %w", hostID, payload.Title, err)
	}

	propertyItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Property])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:properties for host_id=%s title=%s: %w", hostID, payload.Title, err)
	}

	return &propertyItem, nil
}

func (pr *PropertyRepository) GetAllProperties(ctx context.Context) ([]*Property, error) {
	stmt := `
		SELECT
			id, title, sub_title, price, host_id, max_guests, created_at, updated_at
		FROM
			properties
	`
	rows, err := pr.server.DB.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get all properties query: %w", err)
	}
	defer rows.Close()

	properties, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Property])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows from table:properties: %w", err)
	}

	return properties, nil
}

func (pr *PropertyRepository) GetPropertyByID(ctx context.Context, propertyID string) (*Property, error) {
	stmt := `
		SELECT
			id, title, sub_title, price, host_id, max_guests, created_at, updated_at
		FROM
			properties
		WHERE
			id = @id
		`

	rows, err := pr.server.DB.Query(ctx, stmt, pgx.NamedArgs{
		"id": propertyID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute get property by id for property_id %v: %w", propertyID, err)
	}

	propertyItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Property])

	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:properties for property_id=%s: %w", propertyID, err)
	}

	return &propertyItem, nil

}

func (pr *PropertyRepository) UpdateProperty(ctx context.Context, propertyID string, payload *UpdatePropertyPayload) (*Property, error) {
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

	rows, err := pr.server.DB.Query(ctx, stmt, args)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	updatedProperty, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Property])

	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:properties for property_id=%s: %w", propertyID, err)
	}

	return &updatedProperty, nil

}

func (pr *PropertyRepository) DeleteProperty(ctx context.Context, propertyID string) error {
	stmt := `
		DELETE FROM properties
		WHERE id = @id
	`

	result, err := pr.server.DB.Exec(ctx, stmt, pgx.NamedArgs{
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
