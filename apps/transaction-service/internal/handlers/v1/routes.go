package v1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/yarmaq/pkg/middleware"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/handlers/v1/transaction"
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
	transactionRouter := r.router.Group("/transactions")
	transactionRouter.Get("/ping", func(c *fiber.Ctx) error {
		logger.Info("ping received", time.Now())
		return c.Status(200).JSON(fiber.Map{"message": "pong"})
	})
	transactionRouter.Use(
		middleware.AuthMiddleware(r.handler.JWKS),
		middleware.UserContextMiddleware(),
		middleware.RateLimiter(3, time.Second * 10),
	)
	transaction.NewRouter(transactionRouter, r.handler.Transaction, r.handler.JWKS).SetupRoutes()
}
