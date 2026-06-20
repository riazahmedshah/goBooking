package booking

import (
	"context"
	"fmt"

	"github.com/riazahmedshah/go-booking/internal/lib/utils"
	"github.com/riazahmedshah/go-booking/internal/server"
)

type BookingService struct {
	server      *server.Server
	bookingRepo *BookingRepository
}

func NewBookingService(server *server.Server, bookingRepo *BookingRepository) *BookingService {
	return &BookingService{
		server:      server,
		bookingRepo: bookingRepo,
	}
}

func (b *BookingService) CreateBooking(ctx context.Context, userID int, payload *CreateBookingPayload) (any, error) {

	lockKey := fmt.Sprintf("booking:%d", payload.PropertyID)
	lockCtx, cancel, err := b.server.Locker.WithContext(ctx, lockKey)
	if err != nil {
		return nil, err
	}
	defer cancel()

	booking, err := b.bookingRepo.CreateBooking(lockCtx, payload)
	if err != nil {
		return nil, err
	}

	key, err := utils.GenerateIdempotencyKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate idempotency key: %w", err)
	}

	idempotencyData, err := b.bookingRepo.CreateIdempotencyKey(lockCtx, key, booking.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create idempotency key: %w", err)
	}

	return idempotencyData.IdemKey, nil

}

func (b *BookingService) ConfirmBooking(ctx context.Context, key string, payload *ConfirmBookingPayload) (any, error) {
	tx, err := b.server.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	idempotencyData, err := b.bookingRepo.GetIdempotencyKeyWithLock(ctx, tx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get idempotency key: %w", err)
	}

	if idempotencyData.IsFinalized {
		return nil, fmt.Errorf("booking is already finalized")
	}

	booking, err := b.bookingRepo.ConfirmBooking(ctx, tx, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm booking: %w", err)
	}

	if err := b.bookingRepo.FinalizeIdempotencyKey(ctx, tx, key); err != nil {
		return nil, fmt.Errorf("failed to finalize idempotency key: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return booking, nil
}
