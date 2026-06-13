package main

import (
    "log"
    "github.com/gofiber/fiber/v3"
    "github.com/Launchkit-org/LaunchKit/shared/config"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("failed to load config: %v", err)
    }

    app := fiber.New()

    app.Get("/health", func(c fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status":  "ok",
            "service": "gateway",
            "env":     cfg.Gateway.HTTPAddr,
        })
    })

    app.Listen(cfg.Gateway.HTTPAddr)
}