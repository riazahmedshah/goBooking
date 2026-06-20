package booking

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/riazahmedshah/go-booking/internal/server"
)

type BookingRepository struct {
	server *server.Server
}

func NewBookingRepository(server *server.Server) *BookingRepository {
	return &BookingRepository{
		server: server,
	}
}

func (r *BookingRepository) CreateBooking(ctx context.Context, payload *CreateBookingPayload) (*Booking, error) {
	stmt := `
		INSERT INTO bookings (
			user_id, 
			property_id, 
			total_price, 
			status
		)
		VALUES (
			@user_id, 
			@property_id, 
			@total_price, 
			@status
		)
		RETURNING *
	`

	rows, err := r.server.DB.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":     payload.UserID,
		"property_id": payload.PropertyID,
		"total_price": payload.TotalPrice,
		"status":      payload.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create booking query for user_id=%v property_id=%v: %w", payload.UserID, *payload.PropertyID, err)
	}
	defer rows.Close()

	bookingItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Booking])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:bookings for user_id=%v property_id=%v: %w", *payload.UserID, *payload.PropertyID, err)
	}

	return &bookingItem, nil
}

func (r *BookingRepository) CreateIdempotencyKey(ctx context.Context, idemKey *CreateIdempotencyKeyPayload, bookingId int) (*IdempotencyKey, error) {
	stmt := `
		INSERT INTO idempotency_keys (
			idem_key, 
			booking_id, 
			is_finalized
		)
		VALUES (
			@idem_key, 
			@booking_id, 
			false
		)
		RETURNING *
	`

	rows, err := r.server.DB.Query(ctx, stmt, pgx.NamedArgs{
		"idem_key":   idemKey.IdemKey,
		"booking_id": bookingId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create idempotency key query for key=%v booking_id=%v: %w", *idemKey.IdemKey, bookingId, err)
	}
	defer rows.Close()

	idemKeyItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[IdempotencyKey])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:idempotency_keys for key=%v booking_id=%v: %w", *idemKey.IdemKey, bookingId, err)
	}

	return &idemKeyItem, nil
}

func (r *BookingRepository) ConfirmBooking(ctx context.Context, payload *ConfirmBookingPayload) error {
	stmt := `
		UPDATE bookings
		SET status = @status
		WHERE id = @id
	`
	_, err := r.server.DB.Exec(ctx, stmt, pgx.NamedArgs{
		"status": "confirmed",
		"id":     payload.BookingID,
	})
	if err != nil {
		return fmt.Errorf("failed to execute confirm booking query for id=%d: %w", *payload.BookingID, err)
	}
	return nil
}

func (r *BookingRepository) FinalizeIdempotencyKey(ctx context.Context, key string) error {
	stmt := `
		UPDATE idempotency_keys
		SET is_finalized = true
		WHERE idem_key = @key
	`
	_, err := r.server.DB.Exec(ctx, stmt, pgx.NamedArgs{
		"idem_key": key,
	})
	if err != nil {
		return fmt.Errorf("failed to execute finalize idempotency key query for key=%v: %w", key, err)
	}
	return nil
}

func (r *BookingRepository) GetIdempotencyKeyWithLock(ctx context.Context, key string) (*IdempotencyKey, error) {
	stmt := `
		SELECT 
			id, 
			idem_key, 
			booking_id, 
			is_finalized
		FROM 
			idempotency_keys
		WHERE 
			idem_key = @key
		FOR UPDATE
	`
	rows, err := r.server.DB.Query(ctx, stmt, pgx.NamedArgs{
		"idem_key": key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get idempotency key with lock query for key=%s: %w", key, err)
	}
	defer rows.Close()

	idemKeyItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[IdempotencyKey])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:idempotency_keys for key=%v: %w", key, err)
	}

	return &idemKeyItem, nil
}
