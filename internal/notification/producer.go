package notification

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

func EnqueueBookingCompletionTask(client *asynq.Client, task *BookingCompletionTask) error {
	payload, err := json.Marshal(task)

	if err != nil {
		return err
	}

	asynqTask := asynq.NewTask(TaskBookingCompletion, payload)

	_, err = client.Enqueue(asynqTask)
	return err
}
