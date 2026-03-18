package webhook

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	router  fiber.Router
	handler Handler
}

func NewRouter(router fiber.Router, handler Handler) *Router {
	return &Router{
		router: router,
		handler: handler,
	}
}

func (r *Router) SetupRoutes() {
	r.router.Post("/sync", r.handler.HandleZitadelSync)
}