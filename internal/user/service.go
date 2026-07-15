package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/riazahmedshah/go-booking/internal/lib/utils"
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
	// if err == nil && errors.Is(err, pgx.ErrNoRows) {
	// 	slog.Error("database failure during email check", "error", err)
	// 	return ErrInternal
	// }
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

func (us *UserService) Login(ctx context.Context, payload *LoginPayload) (string, error) {
	exixtingUser, err := us.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("invalid email or password")
		}
		slog.Error("database failure during login", "error", err)
		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(exixtingUser.Password), []byte(payload.Password))
	if err != nil {
		return "", errors.New("invalid email/password")
	}

	token, err := utils.GenerateJWTToken(exixtingUser.ID, exixtingUser.Role, []byte(us.server.Config.JWT.SecretKey))
	if err != nil {
		slog.Error("failed to generate JWT token", "error", err)
		return "", ErrInternal
	}

	return token, nil
}
