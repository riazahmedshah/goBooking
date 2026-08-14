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

const (
	msgBookingFailed       = "unable to confirm your booking, please try again"
	msgCreateBookingFailed = "unable to initiate booking request, please try after sometime"
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

	if err != nil && !rueidis.IsRedisNil(err) {
		return nil, errs.New(http.StatusInternalServerError, msgCreateBookingFailed, err)
	}

	detachedCtx := context.WithoutCancel(ctx)

	booking, err := b.bookingRepo.CreateBooking(detachedCtx, userID, payload)
	if err != nil {
		_ = b.server.RedisClient.Do(detachedCtx, b.server.RedisClient.B().Del().Key(holdKey).Build())
		return nil, errs.New(http.StatusInternalServerError, msgCreateBookingFailed, err)
	}

	key, err := utils.GenerateIdempotencyKey()
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgCreateBookingFailed, err)
	}

	idempotencyData, err := b.bookingRepo.CreateIdempotencyKey(detachedCtx, key, booking.ID)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgCreateBookingFailed, err)
	}

	return idempotencyData.Key, nil

}

func (b *BookingService) ConfirmBooking(ctx context.Context, key string, userID string, payload *ConfirmBookingPayload) (any, error) {

	tx, err := b.server.DB.Begin(ctx)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgBookingFailed, err)
	}
	defer tx.Rollback(ctx)

	idempotencyData, err := b.bookingRepo.GetIdempotencyKeyWithLock(ctx, tx, key)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgBookingFailed, err)
	}

	if idempotencyData.IsFinalized {
		return nil, errs.ErrDuplicateBooking
	}

	booking, err := b.bookingRepo.ConfirmBooking(ctx, tx, payload)
	if err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgBookingFailed, err)
	}

	if err := b.bookingRepo.FinalizeIdempotencyKey(ctx, tx, key); err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgBookingFailed, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errs.New(http.StatusInternalServerError, msgBookingFailed, err)
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
