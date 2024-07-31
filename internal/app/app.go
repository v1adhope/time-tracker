package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/configs"
	v1 "github.com/v1adhope/time-tracker/internal/controllers/v1"
	"github.com/v1adhope/time-tracker/internal/usecases"
	"github.com/v1adhope/time-tracker/internal/usecases/repositories"
	"github.com/v1adhope/time-tracker/pkg/httpserver"
	"github.com/v1adhope/time-tracker/pkg/logger"
	"github.com/v1adhope/time-tracker/pkg/postgresql"
)

func Run(cfg *configs.Config, log logger.Logger) error {
	mainCtx := context.Background()

	postgres, err := postgresql.Build(mainCtx, cfg.Postgres, "migrations")
	if err != nil {
		return err
	}
	defer postgres.Close()
	log.Info("postgres driver was succsesfully up")

	if cfg.Postgres.WithMigrate {
		if err := postgres.MigrateUp(); err != nil {
			return err
		}
		log.Info("postgres migration was succeeded")
	}

	gin.SetMode(cfg.Gin.Mode)

	repos := repositories.New(postgres)

	usecases := usecases.New(repos)

	if err := v1.RegisterCustomValidations(); err != nil {
		return err
	}
	log.Info("custom validation rules was connected")

	handler := gin.New()

	v1.Handle(&v1.Router{
		Handler:  handler,
		Usecases: usecases,
		Log:      log,
	})

	httpserver.New(handler, &cfg.Server).Run()

	return nil
}
