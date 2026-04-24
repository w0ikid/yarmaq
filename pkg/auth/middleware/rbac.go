package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func RBACMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles, ok := c.Locals("roles").([]string)
		if !ok || len(roles) == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "No roles assigned"})
		}

		for _, allowed := range allowedRoles {
			for _, role := range roles {
				if role == allowed {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}
}
