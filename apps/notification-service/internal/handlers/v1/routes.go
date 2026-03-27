package v1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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

	// routes
	notificationRouter := r.router.Group("/notifications")
	notificationRouter.Get("/ping", func(c *fiber.Ctx) error {
		logger.Info("ping received", time.Now())
		return c.Status(200).JSON(fiber.Map{"message": "pong"})
	})
}
