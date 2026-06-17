package response

import (
	"github.com/gofiber/fiber/v3"
)

func OK(c fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Created(c fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Accepted(c fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func BadRequest(c fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"data":    data,
		"error":   message,
	})
}

func TooManyRequests(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"error":   message,
	})
}

func InternalServerError(c fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"error":   "An unexpected server error occurred. Please try again later.",
	})
}

func Unauthorized(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"error":   message, 
	})
}

func NotFound(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}

func Conflict(c fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}

func Forbidden(c fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"success": false,
		"data":    data,
		"error":   message,
	})
}