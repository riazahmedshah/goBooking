package property

import "time"

type Property struct {
	ID        string     `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	SubTitle  *string    `json:"subTitle" db:"sub_title"`
	Price     *float64   `json:"price" db:"price"`
	HostID    string     `json:"hostId" db:"host_id"`
	MaxGuests *int       `json:"maxGuest" db:"max_guests"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}
