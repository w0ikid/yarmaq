package internals

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/account"
	"github.com/w0ikid/yarmaq/pkg/errs"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	AccountDomain account.AccountDomain
	Logger        *zap.SugaredLogger
}

type Handler interface {
	UpdateBalance(c *fiber.Ctx) error
}

type handler struct {
	domain account.AccountDomain
	logger *zap.SugaredLogger
}

func NewHandler(deps HandlerDeps) Handler {
	return &handler{
		domain: deps.AccountDomain,
		logger: deps.Logger.Named("internal_handler"),
	}
}

func (h *handler) UpdateBalance(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid account ID"})
	}

	var req models.UpdateBalanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err = h.domain.UpdateBalanceUsecase.Execute(c.Context(), id, req.Amount, req.OperationType, req.ReferenceID)
	if err != nil {
		h.logger.Errorw("failed to update balance", "id", id, "error", err)
		return errs.HandleHTTP(c, err)
	}

	return c.Status(200).Send(nil)
}
