package notification

const TaskBookingCompletion = "email:booking_completion"

type BookingCompletionTask struct {
	BookingID int `json:"booking_id"`
	// PropertyName string    `json:"property_name"`
	// StartDate    time.Time `json:"start_date"`
	// EndDate      time.Time `json:"end_date"`
	// Address      string    `json:"address"`
	// TotalMembers int       `json:"total_members"`
	TotalPrice float64 `json:"total_price"`
}
