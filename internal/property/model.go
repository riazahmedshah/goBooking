package property

import "time"

type Property struct {
	ID        string     `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	SubTitle  *string    `json:"subTitle" db:"sub_title"`
	Price     *float64   `json:"price" db:"price"`
	HostID    string     `json:"hostId" db:"host_id"`
	MaxGuests *int       `json:"maxGuests" db:"max_guests"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}

type PropertyAvailabiliy struct {
	ID          string    `json:"id" db:"id"`
	PropertyID  string    `json:"propertyId" db:"property_id"`
	Date        time.Time `json:"date" db:"date"`
	IsAvailable bool      `json:"isAvailable" db:"is_available"`
	BookingID   *string   `json:"bookingId" db:"booking_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type DayAvailability struct {
	CalendarDate string `json:"calendarDate"`
	Available    bool   `json:"available"`
}

type MonthAvailability struct {
	Month int               `json:"month"`
	Year  int               `json:"year"`
	Days  []DayAvailability `json:"days"`
}
