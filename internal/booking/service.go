package booking

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/rueidis"
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

var ErrPropertyHeld = errors.New("property is currently held by another booking request, please try again later")

func (b *BookingService) CreateBooking(ctx context.Context, userID string, payload *CreateBookingPayload) (any, error) {

	// Example key: "hold:property:prop_123:dates:2026-08-15_2026-08-17"
	holdKey := fmt.Sprintf("hold:property:%s:dates:%s_%s",
		*payload.PropertyID,
		*payload.CheckIn,
		*payload.CheckOut,
	)

	cmd := b.server.RedisClient.B().
		Set().
		Key(holdKey).
		Value(userID). // UserID or Booking Ref as value
		Nx().          // Set Only If Not Exists
		Ex(59 * time.Second).
		Build()

	res := b.server.RedisClient.Do(ctx, cmd)
	// if err := res.Error(); err != nil {
	// 	return nil, fmt.Errorf("failed to process hold in cache: %w", err)
	// }

	isSet, err := res.AsBool()
	if !isSet {
		valCmd := b.server.RedisClient.B().Get().Key(holdKey).Build()
		heldByUserID, getErr := b.server.RedisClient.Do(ctx, valCmd).ToString()

		if getErr == nil && heldByUserID == userID {
			// Same user ne dubara request mari hai
			return nil, fmt.Errorf("your booking is already in progress, please check your payments or wait a moment")
		}
		return nil, ErrPropertyHeld
	}

	// Unexpected actual Redis network/connection error check
	if err != nil && !rueidis.IsRedisNil(err) {
		return nil, fmt.Errorf("failed to process hold in cache: %w", err)
	}

	detachedCtx := context.WithoutCancel(ctx)

	booking, err := b.bookingRepo.CreateBooking(detachedCtx, userID, payload)
	if err != nil {
		_ = b.server.RedisClient.Do(detachedCtx, b.server.RedisClient.B().Del().Key(holdKey).Build())
		return nil, fmt.Errorf("failed to create booking record: %w", err)
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
