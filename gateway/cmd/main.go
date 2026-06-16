package main

import (
	"log"

	"github.com/Launchkit-org/LaunchKit/gateway/internal/app"
	"github.com/Launchkit-org/LaunchKit/shared/config"
	"github.com/Launchkit-org/LaunchKit/shared/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	appLogger := logger.NewLogger(&cfg.Log)

	app, err := app.StartApp(cfg)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("cannot start app")
	}

	if err := Run(app, cfg.Gateway.HTTPAddr); err != nil {
		appLogger.Fatal().Err(err).Msg("cannot run program")
	}
}
