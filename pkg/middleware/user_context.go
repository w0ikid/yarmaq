package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/w0ikid/yarmaq/pkg/ctxkeys"
)

func UserContextMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID, ok := c.Locals("userID").(string)
        if !ok || userID == "" {
            return c.Next()
        }

        roles, _ := c.Locals("roles").([]string)

        ctx := ctxkeys.WithUserContext(c.UserContext(), userID, roles)
        c.SetUserContext(ctx)

        return c.Next()
    }
}
