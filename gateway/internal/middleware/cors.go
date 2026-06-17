package middleware

import (
	"github.com/Launchkit-org/LaunchKit/shared/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func NewCORS(cfg config.CORSConfig) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowHeaders:     cfg.AllowedHeaders,
		AllowMethods:     cfg.AllowedMethods,
		AllowCredentials: true,
	})
}