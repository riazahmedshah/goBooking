package email

func (c *Client) SendConfirmationEmail(to string, bookingID int, totalPrice float64) error {

	data := map[string]any{
		"BookingID": bookingID,
		// "PropertyName": propertyName,
		// "StartDate":    startDate.Format("January 2, 2006"),
		// "EndDate":      endDate.Format("January 2, 2006"),
		// "Address":      address,
		// "TotalMembers": totalMembers,
		"TotalPrice": totalPrice,
	}

	return c.SendEmail(to, "Booking Confirmation", "confirmation", data)
}
