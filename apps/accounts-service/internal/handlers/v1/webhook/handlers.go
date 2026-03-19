package webhook

import (
	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/account"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	AccountDomain account.AccountDomain
	Logger        *zap.SugaredLogger
}

type Handler interface {
	HandleZitadelSync(c *fiber.Ctx) error
}

type handler struct {
	accountDomain account.AccountDomain
	logger        *zap.SugaredLogger
}

func NewHandler(deps HandlerDeps) Handler {
	log := deps.Logger.Named("webhook_handler")
	return &handler{
		accountDomain: deps.AccountDomain,
		logger:        log,
	}
}
