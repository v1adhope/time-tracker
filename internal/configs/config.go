package configs

import (
	"fmt"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/v1adhope/time-tracker/pkg/httpserver"
	"github.com/v1adhope/time-tracker/pkg/postgresql"
)

type Config struct {
	Postgres postgresql.Config
	Server   httpserver.Config
}

func Build() (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		return nil, fmt.Errorf("config load: %w", err)
	}

	var cfg Config

	if err := k.Unmarshal("", &cfg.Postgres); err != nil {
		return nil, fmt.Errorf("config unmarshal: postgres: %w", err)
	}

	if err := k.Unmarshal("", &cfg.Server); err != nil {
		return nil, fmt.Errorf("config unmarshal: server: %w", err)
	}

	return &cfg, nil
}
