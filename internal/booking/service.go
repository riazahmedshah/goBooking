package booking

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/riazahmedshah/go-booking/internal/lib/utils"
	"github.com/riazahmedshah/go-booking/internal/notification"
	"github.com/riazahmedshah/go-booking/internal/server"
)

type BookingService struct {
	server       *server.Server
	bookingRepo  *BookingRepository
	notification *notification.NotificationService
}

func NewBookingService(server *server.Server, bookingRepo *BookingRepository, notification *notification.NotificationService) *BookingService {
	return &BookingService{
		server:       server,
		bookingRepo:  bookingRepo,
		notification: notification,
	}
}

func (b *BookingService) CreateBooking(ctx context.Context, userID string, payload *CreateBookingPayload) (any, error) {

	if b.notification == nil {
		slog.Error("CRITICAL: b.notification is NIL inside BookingService!")
	}

	lockKey := fmt.Sprintf("booking:%s", *payload.PropertyID)
	// lockTimeoutCtx, cancelTimeout := context.WithTimeout(ctx, 500*time.Millisecond)
	// defer cancelTimeout()
	detachedCtx := context.WithoutCancel(ctx)

	tryLockCtx, cancelTryLock := context.WithTimeout(detachedCtx, 5*time.Second)
	defer cancelTryLock()
	_, _, err := b.server.Locker.WithContext(tryLockCtx, lockKey)
	if err != nil {
		return nil, fmt.Errorf("property is currently being booked by another request, please try again: %w", err)
	}

	booking, err := b.bookingRepo.CreateBooking(detachedCtx, userID, payload)
	if err != nil {
		return nil, err
	}

	key, err := utils.GenerateIdempotencyKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate idempotency key: %w", err)
	}

	idempotencyData, err := b.bookingRepo.CreateIdempotencyKey(detachedCtx, key, booking.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create idempotency key: %w", err)
	}

	return idempotencyData.Key, nil

}

func (b *BookingService) ConfirmBooking(ctx context.Context, key string, userID string, payload *ConfirmBookingPayload) (any, error) {

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

	if err := b.notification.EnqueueBookingCompletionTask(&notification.BookingCompletionTask{
		UserID:     userID,
		BookingID:  booking.ID,
		TotalPrice: booking.TotalPrice,
	}); err != nil {
		slog.Error("failed to enqueue booking completion email", "error", err)
	}

	return booking, nil
}
