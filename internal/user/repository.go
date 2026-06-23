package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/riazahmedshah/go-booking/internal/server"
)

type UserRepository struct {
	server *server.Server
}

func NewUserRepository(server *server.Server) *UserRepository {
	return &UserRepository{
		server: server,
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, payload *CreateUserPayload) error {
	stmt := `
		INSERT INTO users(
			first_name, last_name, email, password
		) 
		VALUES(
			@first_name, @last_name, @email, @password
		)
	`
	_, err := ur.server.DB.Exec(ctx, stmt, pgx.NamedArgs{
		"first_name": payload.FirstName,
		"last_name":  payload.LastName,
		"email":      payload.Email,
		"password":   payload.Password,
	})

	if err != nil {
		return fmt.Errorf("failed to execute create user query for email=%s: %w", payload.Email, err)
	}

	return nil
}
