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

	idem, err := b.bookingRepo.CreateIdempotencyKey(lockCtx, key, booking.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create idempotency key: %w", err)
	}

	return idem.IdemKey, nil

}
