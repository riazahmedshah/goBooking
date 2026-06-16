package property

type Property struct {
	ID        string `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	SubTitle  string `json:"subTitle" db:"sub_title"`
	Image     string `json:"image" db:"image"`
	AddressID int    `json:"addressId" db:"address_id"`
	HostID    int    `json:"hostId" db:"host_id"`
	MaxGuests int    `json:"maxGuest" db:"max_guests"`
}
