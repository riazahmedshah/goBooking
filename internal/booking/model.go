package booking

import "time"

type Booking struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"userId" db:"user_id"`
	PropertyID string    `json:"propertyId" db:"property_id"`
	TotalPrice float64   `json:"totalPrice" db:"total_price"`
	Status     string    `json:"status" db:"status"`
	CheckIn    time.Time `json:"checkIn" db:"check_in"`
	CheckOut   time.Time `json:"checkOut" db:"check_out"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type IdempotencyKey struct {
	ID          string    `json:"id" db:"id"`
	Key         string    `json:"key" db:"key"`
	BookingID   string    `json:"bookingId" db:"booking_id"`
	IsFinalized bool      `json:"isFinalized" db:"is_finalized"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}
