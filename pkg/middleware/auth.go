package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "github.com/w0ikid/yarmaq/pkg/jwks"
    "go.uber.org/zap"
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
        c.Locals("claims", claims)
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

func ServiceOnlyMiddleware(serviceName string, logger *zap.SugaredLogger, allowedServices ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        rawClaims := c.Locals("claims")
        claims, ok := rawClaims.(jwt.MapClaims)
        if !ok {
            mapClaims, ok2 := rawClaims.(map[string]interface{})
            if !ok2 {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token context"})
            }
            claims = jwt.MapClaims(mapClaims)
        }

        clientID, _ := claims["client_id"].(string)
        if clientID == "" {
            clientID, _ = claims["azp"].(string)
        }
        if clientID == "" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Service token required"})
        }

        allowed := false
        for _, svc := range allowedServices {
            if svc == clientID {
                allowed = true
                break
            }
        }
        if !allowed {
            logger.Warnw("unauthorized service call", "client_id", clientID, "path", c.Path())
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Service not allowed"})
        }

        logger.Infow("internal service call", "from", clientID, "path", c.Path())
        return c.Next()
    }
}