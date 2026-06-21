package notification

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TaskBookingCompletion = "email:booking_completion"

type BookingCompletionTask struct {
	BookingID    string  `json:"booking_id"`
	UserEmail    string  `json:"user_email"`
	PropertyName string  `json:"property_name"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	Adress       string  `json:"address"`
	TotalMembers int     `json:"total_members"`
	TotalPrice   float64 `json:"total_price"`
}

func EnqueueBookingCompletionTask(client asynq.Client, task *BookingCompletionTask) error {
	payload, err := json.Marshal(task)

	if err != nil {
		return err
	}

	asynqTask := asynq.NewTask(TaskBookingCompletion, payload)

	_, err = client.Enqueue(asynqTask)
	return err
}
