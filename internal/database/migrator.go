package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/url"

	"github.com/riazahmedshah/go-booking/internal/config"

	"github.com/jackc/pgx/v5"
	tern "github.com/jackc/tern/v2/migrate"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(ctx context.Context, cfg *config.Config) error {
	hostPort := net.JoinHostPort(cfg.Database.Host, cfg.Database.Port)

	// URL-encode the password
	encodedPassword := url.QueryEscape(cfg.Database.Password)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		encodedPassword,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	m, err := tern.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		return fmt.Errorf("constructing database migrator: %w", err)
	}
	subtree, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("retrieving database migrations subtree: %w", err)
	}
	if err := m.LoadMigrations(subtree); err != nil {
		return fmt.Errorf("loading database migrations: %w", err)
	}
	from, err := m.GetCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("retreiving current database migration version")
	}
	if err := m.Migrate(ctx); err != nil {
		return err
	}
	if from == int32(len(m.Migrations)) {
		slog.Info("database schema up to date", "version", len(m.Migrations))
	} else {
		slog.Info("migrated database schema", "from", from, "to", len(m.Migrations))
	}
	return nil
}
