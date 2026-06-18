package router

import (
	"github.com/Launchkit-org/LaunchKit/gateway/internal/handler"
	"github.com/Launchkit-org/LaunchKit/gateway/internal/middleware"
	response "github.com/Launchkit-org/LaunchKit/shared/responses"
	"github.com/gofiber/fiber/v3"
)

type RouteDependencies struct {
	AuthHandler    *handler.AuthHandler
	AuthMiddleware middleware.AuthMiddleware
	AuthLimiter    fiber.Handler
	PublicLimiter  fiber.Handler
	ApiLimiter     fiber.Handler
}

func SetUpRoutes(app *fiber.App, deps RouteDependencies) {
	app.Get("/health", deps.PublicLimiter, func(c fiber.Ctx) error {
		return response.OK(c, "health check done", nil)
	})

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Get("/nonce", deps.AuthLimiter, deps.AuthHandler.GetNonce)
	auth.Post("/verify", deps.AuthLimiter, deps.AuthHandler.Verify)
	auth.Post("/logout", deps.AuthMiddleware.Auth, deps.ApiLimiter, deps.AuthHandler.Logout)
}


