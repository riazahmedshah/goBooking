package booking

type Booking struct {
	ID         int     `json:"id" db:"id"`
	UserID     int     `json:"userId" db:"user_id"`
	PropertyID int     `json:"propertyId" db:"property_id"`
	TotalPrice float64 `json:"totalPrice" db:"total_price"`
	Status     string  `json:"status" db:"status"`
}

type IdempotencyKey struct {
	ID          int    `json:"id" db:"id"`
	Key         string `json:"key" db:"key"`
	BookingID   int    `json:"bookingId" db:"booking_id"`
	IsFinalized bool   `json:"isFinalized" db:"is_finalized"`
}
