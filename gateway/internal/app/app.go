package app

import (
	"fmt"

	"github.com/Launchkit-org/LaunchKit/gateway/internal/router"
	"github.com/Launchkit-org/LaunchKit/shared/cache"
	"github.com/Launchkit-org/LaunchKit/shared/config"
	"github.com/Launchkit-org/LaunchKit/shared/logger"
	"github.com/Launchkit-org/LaunchKit/shared/serializer"
	"github.com/Launchkit-org/LaunchKit/shared/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type App struct {
	Server *fiber.App
	Redis  *redis.Client
	Logger zerolog.Logger
}

func StartApp(cfg *config.Config) (*App, error) {

	appLogger := logger.NewLogger(&cfg.Log)

	redis, err := cache.ConnectRedis(&cfg.Redis)
	if err != nil {
		appLogger.Error().Err(err).Msg("connect redis")
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	v := validator.NewValidator()

	app := fiber.New(fiber.Config{

		JSONEncoder:     serializer.Marshal,
		JSONDecoder:     serializer.Unmarshal,
		IdleTimeout:     cfg.Gateway.IdleTimeout,
		ReadTimeout:     cfg.Gateway.ReadTimeout,
		WriteTimeout:    cfg.Gateway.WriteTimeout,
		BodyLimit:       10 * 1024 * 1024,
		StructValidator: v,
	})

	router.SetUpRoutes(app)

	return &App{
		Server: app,
		Redis:  redis,

		Logger: appLogger,
	}, nil

}
