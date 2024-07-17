package main

import (
	"log"

	"github.com/v1adhope/time-tracker/internal/app"
	"github.com/v1adhope/time-tracker/internal/configs"
	"github.com/v1adhope/time-tracker/pkg/logger"
)

func main() {
	cfg, err := configs.Build(".env")
	if err != nil {
		log.Fatal(err)
	}

	appLog := logger.New(cfg.Logger.LogLevel)

	if err := app.Run(cfg, appLog); err != nil {
		log.Fatal(err)
	}
}
