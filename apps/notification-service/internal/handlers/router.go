package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	v1 "github.com/w0ikid/yarmaq/apps/notification-service/internal/handlers/v1"
	"go.uber.org/zap"
)

type Router struct {
	app      *fiber.App
	handlers *Handlers
}

func NewRouter(app *fiber.App, handlers *Handlers) *Router {
	return &Router{
		app:      app,
		handlers: handlers,
	}
}

func (r *Router) SetupRoutes(logger *zap.SugaredLogger) {
	// CORS middleware
	r.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	// API v1 routes
	v1Group := r.app.Group("/api/v1")
	v1.NewRouter(v1Group, r.handlers.V1).SetupRoutes(logger)
}
