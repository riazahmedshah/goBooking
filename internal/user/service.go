package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/riazahmedshah/go-booking/internal/server"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	server   *server.Server
	userRepo *UserRepository
}

var (
	ErrEmailAlreadyExists = errors.New("email address is already registered")
	ErrInternal           = errors.New("an unexpected error occurred")
)

func NewUserService(server *server.Server, ur *UserRepository) *UserService {
	return &UserService{
		server:   server,
		userRepo: ur,
	}
}

func (us *UserService) CreateUser(ctx context.Context, payload *CreateUserPayload) error {
	user, err := us.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("database failure during email check", "error", err)
		return ErrInternal
	}

	if err == nil && user != nil {
		return ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		return ErrInternal
	}

	payload.Password = string(hash)
	return us.userRepo.CreateUser(ctx, payload)
}
