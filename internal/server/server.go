package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
	"github.com/riazahmedshah/go-booking/internal/config"
	"github.com/riazahmedshah/go-booking/internal/database"
)

type Server struct {
	Config     *config.Config
	DB         *pgxpool.Pool
	httpServer *http.Server
	Locker     rueidislock.Locker
}

func New(cfg *config.Config) (*Server, error) {
	db, err := database.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// redisClient, err := rueidis.NewClient(rueidis.ClientOption{
	// 	InitAddress: []string{cfg.Redis.Address},
	// })

	locker, err := rueidislock.NewLocker(rueidislock.LockerOption{
		ClientOption: rueidis.ClientOption{
			InitAddress: []string{cfg.Redis.Address},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis client: %w", err)
	}

	server := &Server{
		Config: cfg,
		DB:     db,
		Locker: locker,
	}

	return server, nil
}

func (s *Server) SetupHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:    ":" + s.Config.Server.Port,
		Handler: handler,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}
	slog.Info("starting server", "port", s.Config.Server.Port)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	s.DB.Close()
	s.Locker.Close()
	return nil
}
