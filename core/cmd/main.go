package main

import (
    "github.com/gofiber/fiber/v3"
)

func main() {
    app := fiber.New()

    app.Get("/health", func(c fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status":  "ok",
            "service": "core",
        })
    })

    app.Listen(":8081")
}