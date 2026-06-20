package booking

import "github.com/riazahmedshah/go-booking/internal/server"

type BookingRepository struct {
	server *server.Server
}

func NewBookingRepository(server *server.Server) *BookingRepository {
	return &BookingRepository{
		server: server,
	}
}

func (r *BookingRepository) CreateBooking(booking *CreateBookingPayload) error {
	// Implement the logic to create a booking in the database
	return nil
}

func (r *BookingRepository) CreateIdempotencyKey(key *CreateIdempotencyKeyPayload, bookingId int) error {
	// Implement the logic to create an idempotency key in the database
	return nil
}

func (r *BookingRepository) GetIdempotencyKeyWithLock(key string) (*IdempotencyKey, error) {
	// Implement the logic to retrieve an idempotency key by its key from the database
	return nil, nil
}

func (r *BookingRepository) GetBookingByID(id string) (*Booking, error) {
	// Implement the logic to retrieve a booking by its ID from the database
	return nil, nil
}

func (r *BookingRepository) ConfirmBooking(booking *ConfirmBookingPayload) error {
	// Implement the logic to update a booking in the database
	return nil
}
