package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/middleware"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/users"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/webhook"
	"go.uber.org/zap"
	"time"
)

type Router struct {
	router  fiber.Router
	handler *Handlers
}

func NewRouter(router fiber.Router, handler *Handlers) *Router {
	return &Router{
		router:  router,
		handler: handler,
	}
}

func (r *Router) SetupRoutes(logger *zap.SugaredLogger) {
	r.router.Get("/ping", func(c *fiber.Ctx) error {
		logger.Info("ping received", time.Now())
		return c.Status(200).JSON(fiber.Map{"message": "pong"})
	})

	// routes
	usersRouter := r.router.Group("/users")
	usersRouter.Use(middleware.AuthMiddleware(r.handler.JWKS))
	users.NewRouter(usersRouter, r.handler.Users).SetupRoutes()

	webhookRouter := r.router.Group("/webhook")
	webhook.NewRouter(webhookRouter, r.handler.Webhook).SetupRoutes()
}
