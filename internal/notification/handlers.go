package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

func (n *NotificationService) handleBookingCompletion(ctx context.Context, t *asynq.Task) error {
	var p BookingCompletionTask

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal weekly report email payload: %w", err)
	}

	slog.Info("Sending booking completion email", "booking_id", p.BookingID)

	// TODO: Get User email from AuthService.
	userEmail := "riyazsh360@gmail.com"
	// TODO: Integrate SMTP client.

	if err := n.emailClient.SendConfirmationEmail(
		userEmail,
		p.BookingID,
		// p.PropertyName,
		// p.StartDate,
		// p.EndDate,
		// p.Address,
		// p.TotalMembers,
		p.TotalPrice,
	); err != nil {
		return fmt.Errorf("failed to send booking completion email: %w", err)
	}

	slog.Info("Booking completion email sent successfully", "booking_id", p.BookingID, "userEmail", userEmail)
	return nil
}
