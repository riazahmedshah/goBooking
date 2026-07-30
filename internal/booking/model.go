package booking

import "time"

type Booking struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"userId" db:"user_id"`
	PropertyID int       `json:"propertyId" db:"property_id"`
	TotalPrice float64   `json:"totalPrice" db:"total_price"`
	Status     string    `json:"status" db:"status"`
	CheckIn    time.Time `json:"checkIn" db:"check_in"`
	CheckOut   time.Time `json:"checkOut" db:"check_out"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type IdempotencyKey struct {
	ID          int    `json:"id" db:"id"`
	IdemKey     string `json:"idemKey" db:"idem_key"`
	BookingID   int    `json:"bookingId" db:"booking_id"`
	IsFinalized bool   `json:"isFinalized" db:"is_finalized"`
}
