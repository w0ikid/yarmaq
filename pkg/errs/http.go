package errs

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func HandleHTTP(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrValidation):
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, ErrUnauthorized):
		return c.Status(403).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, ErrNotFound):
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
}
