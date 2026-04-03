package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/account"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/ledger"
	"github.com/w0ikid/yarmaq/pkg/ctxkeys"
	"github.com/w0ikid/yarmaq/pkg/errs"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	AccountDomain account.AccountDomain
	LedgerDomain  ledger.LedgerDomain
	Logger        *zap.SugaredLogger
}

type Handler interface {
	CreateAccount(c *fiber.Ctx) error
	GetMyAccounts(c *fiber.Ctx) error
	GetAccountsByUserID(c *fiber.Ctx) error
	GetAccountByNumberAndCurrency(c *fiber.Ctx) error
	GetAccountByUserIDAndCurrency(c *fiber.Ctx) error

	// ledger handlers
	GetByAccountID(c *fiber.Ctx) error
}

type handler struct {
	accountDomain account.AccountDomain
	ledgerDomain  ledger.LedgerDomain
	logger        *zap.SugaredLogger
}

func NewHandler(deps HandlerDeps) Handler {
	return &handler{
		accountDomain: deps.AccountDomain,
		ledgerDomain:  deps.LedgerDomain,
		logger:        deps.Logger.Named("account_handler"),
	}
}

func (h *handler) CreateAccount(c *fiber.Ctx) error {
	var req CreateAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := ctxkeys.GetUserID(c.UserContext())
	acc := models.Account{
		UserID:   &userID,
		Type:     models.AccountTypeUser,
		Currency: req.Currency,
	}

	created, err := h.accountDomain.CreateUsecase.Execute(c.Context(), acc)
	if err != nil {
		h.logger.Errorw("failed to create account", "error", err)
		return errs.HandleHTTP(c, err)
	}

	return c.Status(201).JSON(created)
}

func (h *handler) GetMyAccounts(c *fiber.Ctx) error {
	accounts, err := h.accountDomain.GetAccountUsecase.ExecuteAllByUserID(c.Context(), ctxkeys.GetUserID(c.UserContext()))
	if err != nil {
		h.logger.Errorw("failed to get user accounts", "error", err)
		return errs.HandleHTTP(c, err)
	}

	return c.Status(200).JSON(accounts)
}

func (h *handler) GetAccountsByUserID(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "user_id is required"})
	}

	accounts, err := h.accountDomain.GetAccountUsecase.ExecuteAllByUserID(c.Context(), userID)
	if err != nil {
		h.logger.Errorw("failed to get accounts by user id", "user_id", userID, "error", err)
		return errs.HandleHTTP(c, err)
	}

	return c.Status(200).JSON(accounts)
}

func (h *handler) GetAccountByNumberAndCurrency(c *fiber.Ctx) error {
	number := c.Query("number")
	currency := c.Query("currency")
	if number == "" || currency == "" {
		return c.Status(400).JSON(fiber.Map{"error": "number and currency are required"})
	}

	acc, err := h.accountDomain.GetAccountUsecase.ExecuteByNumberAndCurrency(c.Context(), number, currency)
	if err != nil {
		h.logger.Errorw("failed to get account by number and currency", "number", number, "currency", currency, "error", err)
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

	acc, err := h.accountDomain.GetAccountUsecase.ExecuteByUserIDAndCurrency(c.Context(), userID, currency)
	if err != nil {
		h.logger.Errorw("failed to get account by user id and currency", "user_id", userID, "currency", currency, "error", err)
		return errs.HandleHTTP(c, err)
	}
	if acc == nil {
		return errs.HandleHTTP(c, errs.ErrNotFound)
	}

	return c.Status(200).JSON(acc)
}
