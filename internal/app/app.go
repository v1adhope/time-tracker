package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/configs"
	v1 "github.com/v1adhope/time-tracker/internal/controllers/v1"
	"github.com/v1adhope/time-tracker/internal/usecases"
	"github.com/v1adhope/time-tracker/internal/usecases/repositories"
	"github.com/v1adhope/time-tracker/pkg/httpserver"
	"github.com/v1adhope/time-tracker/pkg/postgresql"
)

func Run() error {
	mainCtx := context.Background()

	cfg, err := configs.Build()
	if err != nil {
		return err
	}

	postgres, err := postgresql.Build(mainCtx, cfg.Postgres)
	if err != nil {
		return err
	}
	defer postgres.Close()

	if cfg.Postgres.WithMigrate {
		if err := postgres.Migrate(); err != nil {
			return err
		}
	}

	repos := repositories.New(postgres)

	usecases := usecases.New(repos)

	handler := gin.New()

	v1.Handle(&v1.Router{
		Handler:  handler,
		Usecases: usecases,
	})

	httpserver.New(handler, &cfg.Server).Run()

	return nil
}
