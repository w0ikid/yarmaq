package transaction

import (
	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/yarmaq/pkg/jwks"
)

type Router struct {
	router  fiber.Router
	handler Handler
	jwks    *jwks.JWKS
}

func NewRouter(router fiber.Router, handler Handler, j *jwks.JWKS) *Router {
	return &Router{
		router:  router,
		handler: handler,
		jwks:    j,
	}
}

func (r *Router) SetupRoutes() {
	r.router.Post("/", r.handler.CreateTransaction)
	r.router.Post("/transfer", r.handler.CreateTransfer)
	r.router.Post("/deposit", r.handler.CreateDeposit)
	r.router.Post("/withdraw", r.handler.CreateWithdraw)
	r.router.Post("/exchange", r.handler.CreateExchange)
	r.router.Get("/:id", r.handler.GetTransaction)
}
