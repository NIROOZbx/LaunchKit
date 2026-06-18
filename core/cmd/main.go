package main

import (
	"log"

	"github.com/Launchkit-org/LaunchKit/core/internal/app"
	"github.com/Launchkit-org/LaunchKit/shared/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	a, err := app.StartApp(cfg)
	if err != nil {
		log.Fatalf("cannot start app: %v", err)
	}

	if err := Run(a); err != nil {
		a.Logger.Fatal().Err(err).Msg("core service exited with error")
	}
}