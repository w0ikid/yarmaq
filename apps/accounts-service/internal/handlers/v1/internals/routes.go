package internals

import (
	"github.com/gofiber/fiber/v2"
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
	r.router.Post("/accounts/:id/balance", r.handler.UpdateBalance)
	r.router.Get("/accounts/by-number", r.handler.GetAccountByNumberAndCurrency)
	r.router.Get("/accounts/by-user-currency", r.handler.GetAccountByUserIDAndCurrency)
	r.router.Get("/accounts/by-type-currency", r.handler.GetAccountByTypeAndCurrency)
}
