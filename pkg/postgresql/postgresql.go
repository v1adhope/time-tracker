package postgresql

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User        string `koanf:"APP_POSTGRES_USER"`
	Password    string `koanf:"APP_POSTGRES_PASSWORD"`
	Host        string `koanf:"APP_POSTGRES_HOST"`
	Port        string `koanf:"APP_POSTGRES_PORT"`
	DBName      string `koanf:"APP_POSTGRES_DB_NAME"`
	Query       string `koanf:"APP_POSTGRES_QUERY"`
	WithMigrate bool   `koanf:"APP_POSTGRES_WITH_MIGRATION"`
}

type Postgres struct {
	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func Build(ctx context.Context, cfg Config) (*Postgres, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Query,
	)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("postgresql: new: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgresql: ping: %w", err)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &Postgres{pool, builder}, nil
}

func (p *Postgres) Close() {
	p.Pool.Close()
}

func (p *Postgres) Migrate() error {
	m, err := migrate.New("file://migrations", p.Pool.Config().ConnString())
	if err != nil {
		return fmt.Errorf("postgresql: migrate: new: %w", err)
	}
	defer m.Close()

	if err := m.Down(); err != migrate.ErrNoChange && err != nil {
		return fmt.Errorf("postgresql: migrate: down: %w", err)
	}

	if err := m.Up(); err != migrate.ErrNoChange && err != nil {
		return fmt.Errorf("postgresql: migrate: up: %w", err)
	}

	return nil
}
