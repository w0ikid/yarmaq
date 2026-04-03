package v1

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/yarmaq/pkg/middleware"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/account"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/internals"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/ledger"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/v1/webhook"
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
	accountsRouter := r.router.Group("/accounts")
	accountsRouter.Get("/ping", func(c *fiber.Ctx) error {
		logger.Info("ping received", time.Now())
		return c.Status(200).JSON(fiber.Map{"message": "pong"})
	})

	accountsRouter.Use(
		middleware.AuthMiddleware(r.handler.JWKS),
		middleware.UserContextMiddleware(),
	)
	account.NewRouter(accountsRouter, r.handler.Account).SetupRoutes()

	internalRouter := r.router.Group("/internal")
	internalRouter.Use(
		middleware.AuthMiddleware(r.handler.JWKS),
		middleware.ServiceOnlyMiddleware("accounts-service", logger, "transaction-service"),
	)
	internals.NewRouter(internalRouter, r.handler.Internal).SetupRoutes()

	ledgerRouter := r.router.Group("/ledger")
	ledgerRouter.Use(middleware.AuthMiddleware(r.handler.JWKS))
	ledger.NewRouter(ledgerRouter, r.handler.Ledger).SetupRoutes()

	webhookRouter := r.router.Group("/webhook")
	webhook.NewRouter(webhookRouter, r.handler.Webhook).SetupRoutes()
}
