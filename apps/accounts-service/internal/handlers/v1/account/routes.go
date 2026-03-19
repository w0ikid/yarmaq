package account

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
	r.router.Post("/", r.handler.CreateAccount)
	r.router.Get("/:id", r.handler.GetAccount)
	r.router.Post("/:id/balance", r.handler.UpdateBalance)
}
