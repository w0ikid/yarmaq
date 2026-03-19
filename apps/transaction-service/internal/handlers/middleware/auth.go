package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/w0ikid/yarmaq/pkg/jwks"
)

func AuthMiddleware(j *jwks.JWKS) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
        }

        claims, err := j.Validate(authHeader)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
        }

        if sub, ok := claims["sub"].(string); ok {
            c.Locals("userID", sub)
        }
		
        if rolesRaw, ok := claims["urn:zitadel:iam:org:project:roles"].(map[string]interface{}); ok {
            roles := make([]string, 0, len(rolesRaw))
            for role := range rolesRaw {
                roles = append(roles, role)
            }
            c.Locals("roles", roles)
        }

        return c.Next()
    }
}