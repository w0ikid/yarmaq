package webhook

import (
	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/users"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	CreateUser users.CreateUsecase
	UpdateUser users.UpdateUserUsecase
	Logger     *zap.SugaredLogger
}

type Handler interface {
	HandleZitadelSync(c *fiber.Ctx) error
}

type handler struct {
	createUser users.CreateUsecase
	updateUser users.UpdateUserUsecase
	logger     *zap.SugaredLogger
}

func NewHandler(deps HandlerDeps) Handler {
	log := deps.Logger.Named("webhook_handler")
	return &handler{
		createUser: deps.CreateUser,
		updateUser: deps.UpdateUser,
		logger:     log,
	}
}
