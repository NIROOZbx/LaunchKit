package router

import (
	response "github.com/Launchkit-org/LaunchKit/shared/responses"
	"github.com/gofiber/fiber/v3"
)

func SetUpRoutes(app *fiber.App) {
	app.Get("health", func(c fiber.Ctx) error {
		return response.OK(c, "health check done", nil)
	})
}
