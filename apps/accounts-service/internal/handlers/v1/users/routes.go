package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/handlers/middleware"
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
	r.router.Post("/", r.handler.CreateUser)
	r.router.Get("/", r.handler.GetUser)
	r.router.Get("/protected", middleware.RBACMiddleware("manager", "hr"), r.handler.ProtectedUserEndpoint)
}