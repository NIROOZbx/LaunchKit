package app

import (
	"fmt"
	"time"

	"github.com/Launchkit-org/LaunchKit/gateway/internal/client"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/constants"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/handler"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/middleware"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/router"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/service"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/sessions"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/store"
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
	Server     *fiber.App
	Redis      *redis.Client
	Logger     zerolog.Logger
	CoreClient *client.CoreClient
}

func StartApp(cfg *config.Config) (*App, error) {

	appLogger := logger.NewLogger(&cfg.Log)

	redis, err := cache.ConnectRedis(&cfg.Redis)
	if err != nil {
		appLogger.Error().Err(err).Msg("connect redis")
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	// 1. Initialize Clients
	coreClient, err := client.NewCoreClient(cfg.Core.GRPCAddr)
	if err != nil {
		appLogger.Error().Err(err).Msg("failed to initialize core gRPC client")
		return nil, fmt.Errorf("failed to initialize core gRPC client: %w", err)
	}

	// 2. Initialize Stores
	sessionStore := sessions.NewStore(redis)
	c := cache.NewRedisCache(redis)
	authStore := store.NewRedisAuthStore(c)

	// 3. Initialize Services
	authService := service.NewAuthService(authStore, sessionStore, coreClient, &cfg.Jwt)

	// 4. Initialize Handlers & Middlewares
	authHandler := handler.NewAuthHandler(authService, appLogger, &cfg.Jwt)
	authMid := middleware.NewMiddleware(sessionStore, &cfg.Jwt, appLogger, coreClient)

	// 5. Initialize Rate Limiters
	limiters := initRateLimiters(redis, &cfg.RateLimit)

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

	app.Use(middleware.NewCORS(cfg.Cors))

	router.SetUpRoutes(app, router.RouteDependencies{
		AuthHandler:    authHandler,
		AuthMiddleware: authMid,
		AuthLimiter:    limiters.Auth,
		PublicLimiter:  limiters.Public,
		ApiLimiter:     limiters.API,
	})

	return &App{
		Server:     app,
		Redis:      redis,
		Logger:     appLogger,
		CoreClient: coreClient,
	}, nil

}

type Limiters struct {
	Auth   fiber.Handler
	Public fiber.Handler
	API    fiber.Handler
}

func initRateLimiters(redisClient *redis.Client, cfg *config.RateLimitConfig) Limiters {
	return Limiters{
		Auth: middleware.NewRateLimiter(middleware.RateLimiterConfig{
			Redis:      redisClient,
			Limit:      cfg.AuthRequestsPerMinute,
			Window:     time.Minute,
			KeyPrefix:  constants.RateLimitKeyPrefixAuth,
			Identifier: middleware.IdentifierFromIP,
		}),
		Public: middleware.NewRateLimiter(middleware.RateLimiterConfig{
			Redis:      redisClient,
			Limit:      cfg.PublicRequestsPerMinute,
			Window:     time.Minute,
			KeyPrefix:  constants.RateLimitKeyPrefixPublic,
			Identifier: middleware.IdentifierFromIP,
		}),
		API: middleware.NewRateLimiter(middleware.RateLimiterConfig{
			Redis:      redisClient,
			Limit:      cfg.RequestsPerMinute,
			Window:     time.Minute,
			KeyPrefix:  constants.RateLimitKeyPrefixAPI,
			Identifier: middleware.IdentifierFromUser,
		}),
	}
}
