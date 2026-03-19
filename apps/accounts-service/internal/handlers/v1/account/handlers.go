package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/account"
	"github.com/w0ikid/yarmaq/pkg/ctxkeys"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	AccountDomain account.AccountDomain
	Logger        *zap.SugaredLogger
}

type Handler interface {
	CreateAccount(c *fiber.Ctx) error
	GetAccount(c *fiber.Ctx) error
}

type handler struct {
	domain account.AccountDomain
	logger *zap.SugaredLogger
}

func NewHandler(deps HandlerDeps) Handler {
	return &handler{
		domain: deps.AccountDomain,
		logger: deps.Logger.Named("account_handler"),
	}
}

func (h *handler) CreateAccount(c *fiber.Ctx) error {
	var req CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	acc := models.Account{
		UserID:   ctxkeys.GetUserID(c.UserContext()),
		Currency: req.Currency,
	}

	created, err := h.domain.CreateUsecase.Execute(c.Context(), acc)
	if err != nil {
		h.logger.Errorw("failed to create account", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(201).JSON(created)
}

func (h *handler) GetAccount(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid account ID"})
	}

	acc, err := h.domain.GetAccountUsecase.ExecuteByID(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}
	if acc == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Account not found"})
	}

	return c.Status(200).JSON(acc)
}
