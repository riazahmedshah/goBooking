package property

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/riazahmedshah/go-booking/internal/errs"
	"github.com/riazahmedshah/go-booking/internal/server"
)

type PropertyRepository struct {
	server *server.Server
}

func NewPropertyRepository(server *server.Server) *PropertyRepository {
	return &PropertyRepository{server: server}
}

func (pr *PropertyRepository) Createproperty(ctx context.Context, tx pgx.Tx, hostID string, payload *CreatePropertyPayload) (*Property, error) {
	stmt := `
		INSERT INTO properties(
			host_id, title, sub_title, max_guests, price, images
		)
		VALUES (
			@host_id, @title, @sub_title, @max_guests, @price, @images
		)
		RETURNING *
	`

	rows, err := tx.Query(ctx, stmt, pgx.NamedArgs{
		"host_id":    hostID,
		"title":      payload.Title,
		"sub_title":  payload.SubTitle,
		"max_guests": payload.MaxGuests,
		"price":      payload.Price,
		"images":     payload.ImageURLs,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique key violation
			return nil, errs.ErrPropertyTitleExists
		}

		return nil, fmt.Errorf("failed to execute create property query for host_id=%s title=%s: %w", hostID, payload.Title, err)
	}

	propertyItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Property])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:properties for host_id=%s title=%s: %w", hostID, payload.Title, err)
	}

	return &propertyItem, nil
}

func (pr *PropertyRepository) CreateAddress(ctx context.Context, tx pgx.Tx, payload *CreateAddressPayload) (*Address, error) {
	stmt := `
		INSERT INTO addresses (country, state, pincode, city, area, property_id)
		VALUES (@country, @state, @pincode, @city, @area, @property_id)
		RETURNING *
	`

	args := pgx.NamedArgs{
		"country":     payload.Country,
		"state":       payload.State,
		"pincode":     payload.Pincode,
		"city":        payload.City,
		"area":        payload.Area,
		"property_id": payload.PropertyID,
	}

	rows, err := tx.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create address query: %w", err)
	}
	defer rows.Close()

	address, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Address])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:addresses: %w", err)
	}

	return &address, nil
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrPropertyNotFound
		}
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
		return nil, errs.ErrBadUpdateRequest
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
		return errs.ErrPropertyNotFound
	}

	return nil
}

func (pr *PropertyRepository) GetPropertyAvailability(ctx context.Context, propertyID string) ([]*PropertyAvailabiliy, error) {
	stmt := `
		SELECT
			id, property_id, date, is_available, booking_id, created_at, updated_at
		FROM
			property_availability
		WHERE
			property_id = @property_id
	`
	rows, err := pr.server.DB.Query(ctx, stmt, pgx.NamedArgs{
		"property_id": propertyID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	availability, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[PropertyAvailabiliy])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:property_availability for property_id=%s: %w", propertyID, err)
	}

	return availability, nil
}
