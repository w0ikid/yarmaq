package ledger

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/ledger"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	LedgerDomain ledger.LedgerDomain
	Logger       *zap.SugaredLogger
}

type Handler interface {
	GetByAccountID(c *fiber.Ctx) error
}

type handler struct {
	domain ledger.LedgerDomain
	logger *zap.SugaredLogger
}

func NewHandler(deps HandlerDeps) Handler {
	return &handler{
		domain: deps.LedgerDomain,
		logger: deps.Logger.Named("ledger_handler"),
	}
}

func (h *handler) GetByAccountID(c *fiber.Ctx) error {
	idStr := c.Params("accountId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid account ID"})
	}

	entries, err := h.domain.GetLedgerUsecase.ExecuteByAccountID(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(200).JSON(entries)
}
