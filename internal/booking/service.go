package booking

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/rueidis"
	"github.com/riazahmedshah/go-booking/internal/errs"
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

	isSet, err := res.AsBool()
	if !isSet {
		valCmd := b.server.RedisClient.B().Get().Key(holdKey).Build()
		heldByUserID, getErr := b.server.RedisClient.Do(ctx, valCmd).ToString()

		if getErr == nil && heldByUserID == userID {
			return nil, errs.ErrBookingInProgress
		}
		return nil, errs.ErrPropertyHeld
	}

	// Unexpected actual Redis network/connection error check
	if err != nil && !rueidis.IsRedisNil(err) {
		// slog.Error("redis hold error", "err", err, "user_id", userID, "hold_key", holdKey)
		return nil, errs.New(http.StatusInternalServerError, "failed to process booking hold", err)
	}

	detachedCtx := context.WithoutCancel(ctx)

	booking, err := b.bookingRepo.CreateBooking(detachedCtx, userID, payload)
	if err != nil {
		_ = b.server.RedisClient.Do(detachedCtx, b.server.RedisClient.B().Del().Key(holdKey).Build())
		// slog.Error("failed to create booking record in db", "err", err, "user_id", userID)
		return nil, errs.New(http.StatusInternalServerError, "failed to create booking record", err)
	}

	key, err := utils.GenerateIdempotencyKey()
	if err != nil {
		// slog.Error("failed to generate idempotency key", "err", err)
		return nil, errs.New(http.StatusInternalServerError, "failed to process booking security key", err)
	}

	idempotencyData, err := b.bookingRepo.CreateIdempotencyKey(detachedCtx, key, booking.ID)
	if err != nil {
		// slog.Error("failed to store idempotency key in db", "err", err, "booking_id", booking.ID)
		return nil, errs.New(http.StatusInternalServerError, "failed to finalize booking reference", err)
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
