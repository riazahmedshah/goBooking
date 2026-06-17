package property

import "time"

type Property struct {
	ID        int        `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	SubTitle  *string    `json:"subTitle" db:"sub_title"`
	Image     *string    `json:"image" db:"image"`
	AddressID int        `json:"addressId" db:"address_id"`
	HostID    int        `json:"hostId" db:"host_id"`
	MaxGuests *int       `json:"maxGuest" db:"max_guests"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time `json:"updatedAt" db:"updated_at"`
}
