package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/yarmaq/pkg/middleware"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type Router struct {
	router  fiber.Router
	handler Handler
}

func NewRouter(router fiber.Router, handler Handler) *Router {
	return &Router{
		router:  router,
		handler: handler,
	}
}

func (r *Router) SetupRoutes() {
	r.router.Post("/", r.handler.CreateAccount)
	r.router.Get("/me", r.handler.GetMyAccounts)
	r.router.Get("/by-user", middleware.RBACMiddleware(models.RoleSupport, models.RoleAdmin), r.handler.GetAccountsByUserID)
	r.router.Get("/by-number", middleware.RBACMiddleware(models.RoleSupport, models.RoleAdmin), r.handler.GetAccountByNumberAndCurrency)
	r.router.Get("/by-user-currency", middleware.RBACMiddleware(models.RoleSupport, models.RoleAdmin), r.handler.GetAccountByUserIDAndCurrency)

	// ledger routes
	r.router.Get("/:accountId/ledger", middleware.RBACMiddleware(models.RoleSupport), r.handler.GetByAccountID)
}
