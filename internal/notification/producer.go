package notification

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

func (n *NotificationService) EnqueueBookingCompletionTask(task *BookingCompletionTask) error {
	payload, err := json.Marshal(task)

	if err != nil {
		return err
	}

	asynqTask := asynq.NewTask(TaskBookingCompletion, payload)

	_, err = n.client.Enqueue(asynqTask)
	return err
}
