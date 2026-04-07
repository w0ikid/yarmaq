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
	GetAccountByNumberAndCurrency(c *fiber.Ctx) error
	GetAccountByUserIDAndCurrency(c *fiber.Ctx) error
	GetAccountByTypeAndCurrency(c *fiber.Ctx) error
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

func (h *handler) GetAccountByNumberAndCurrency(c *fiber.Ctx) error {
	number := c.Query("number")
	currency := c.Query("currency")
	if number == "" {
		return c.Status(400).JSON(fiber.Map{"error": "number is required"})
	}

	var acc *models.Account
	var err error

	acc, err = h.domain.GetAccountUsecase.ExecuteByNumber(c.Context(), number)

	if err != nil {
		h.logger.Errorw("failed to get account by number", "number", number, "currency", currency, "error", err)
		return errs.HandleHTTP(c, err)
	}
	if acc == nil {
		return errs.HandleHTTP(c, errs.ErrNotFound)
	}

	return c.Status(200).JSON(acc)
}

func (h *handler) GetAccountByUserIDAndCurrency(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	currency := c.Query("currency")
	if userID == "" || currency == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user_id and currency are required"})
	}

	acc, err := h.domain.GetAccountUsecase.ExecuteByUserIDAndCurrency(c.Context(), userID, currency)
	if err != nil {
		h.logger.Errorw("failed to get account by user id and currency", "user_id", userID, "currency", currency, "error", err)
		return errs.HandleHTTP(c, err)
	}
	if acc == nil {
		return errs.HandleHTTP(c, errs.ErrNotFound)
	}

	return c.Status(200).JSON(acc)
}

func (h *handler) GetAccountByTypeAndCurrency(c *fiber.Ctx) error {
	accountType := c.Query("type")
	currency := c.Query("currency")
	if accountType == "" || currency == "" {
		return c.Status(400).JSON(fiber.Map{"error": "type and currency are required"})
	}

	acc, err := h.domain.GetAccountUsecase.ExecuteByTypeAndCurrency(c.Context(), accountType, currency)
	if err != nil {
		h.logger.Errorw("failed to get account by type and currency", "type", accountType, "currency", currency, "error", err)
		return errs.HandleHTTP(c, err)
	}
	if acc == nil {
		return errs.HandleHTTP(c, errs.ErrNotFound)
	}

	return c.Status(200).JSON(acc)
}
