package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

func (n *NotificationService) InitHandlers() {
	// email client initialization logic here
	slog.Info("Initializing email client...")

}

func (n *NotificationService) handleBookingCompletion(ctx context.Context, t *asynq.Task) error {
	var p BookingCompletionTask

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal weekly report email payload: %w", err)
	}

	slog.Info("Sending booking completion email", "booking_id", p.BookingID, "user_email", p.UserEmail)

	// TODO: Get User email from AuthService.
	// TODO: Integrate SMTP client.

	slog.Info("Booking completion email sent successfully", "booking_id", p.BookingID)
	return nil
}
