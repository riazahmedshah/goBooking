package notification

import (
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/riazahmedshah/go-booking/internal/config"
)

type NotificationService struct {
	client *asynq.Client
	server *asynq.Server
}

func NewNotificationServeice(cfg config.Config) *NotificationService {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
	})

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.Address, Password: cfg.Redis.Password, DB: 0},
		asynq.Config{
			Concurrency: 10,
		},
	)
	return &NotificationService{
		client: client,
		server: server,
	}
}

func (n *NotificationService) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskBookingCompletion, n.handleBookingCompletion)

	slog.Info("Starting background workers...")
	if err := n.server.Start(mux); err != nil {
		return err
	}
	return nil
}

func (n *NotificationService) Stop() {
	slog.Info("Shutting down background workers...")
	n.server.Shutdown()
	n.client.Close()
}
