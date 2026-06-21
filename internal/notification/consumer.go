package notification

import (
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/riazahmedshah/go-booking/internal/config"
	"github.com/riazahmedshah/go-booking/internal/lib/email"
)

type NotificationService struct {
	client      *asynq.Client
	server      *asynq.Server
	emailClient *email.Client
}

func NewNotificationService(cfg *config.Config) *NotificationService {
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
		client:      client,
		server:      server,
		emailClient: email.NewClient(cfg.Integration),
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
