package app

import (
	"context"

	"github.com/v1adhope/time-tracker/internal/configs"
	"github.com/v1adhope/time-tracker/pkg/postgresql"
)

func Run() error {
	mainCtx := context.Background()

	cfg, err := configs.Build()
	if err != nil {
		return err
	}

	p, err := postgresql.Build(mainCtx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer p.Close()

	// TODO: think
	if true {
		if err := p.Migrate(); err != nil {
			return err
		}
	}

	return nil
}
