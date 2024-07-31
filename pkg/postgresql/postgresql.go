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
	Pool      *pgxpool.Pool
	Builder   squirrel.StatementBuilderType
	Migration *migrate.Migrate
}

func Build(ctx context.Context, cfg Config, migrationPath string) (*Postgres, error) {
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
		return nil, fmt.Errorf("postgresql: pool: new: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgresql: pool: ping: %w", err)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	migrate, err := migrate.New(fmt.Sprintf("file://%s", migrationPath), connString)
	if err != nil {
		return nil, fmt.Errorf("postgresql: migrate: new: %w", err)
	}

	return &Postgres{pool, builder, migrate}, nil
}

func (p *Postgres) Close() {
	p.Pool.Close()
	p.Migration.Close()
}

func (p *Postgres) MigrateUp() error {
	if err := p.Migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("postgresql: migrate: up: %w", err)
	}

	return nil
}

func (p *Postgres) MigrateDown() error {
	if err := p.Migration.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("postgresql: migrate: down: %w", err)
	}

	return nil
}
