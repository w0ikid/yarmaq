package transaction

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/usecase/transaction"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type HandlerDeps struct {
	TransactionDomain transaction.TransactionDomain
	Logger            *zap.SugaredLogger
}

type Handler interface {
	CreateTransaction(c *fiber.Ctx) error
	GetTransaction(c *fiber.Ctx) error
}

type handler struct {
	domain transaction.TransactionDomain
	logger *zap.SugaredLogger
}

func NewHandler(deps HandlerDeps) Handler {
	return &handler{
		domain: deps.TransactionDomain,
		logger: deps.Logger.Named("transaction_handler"),
	}
}

func (h *handler) CreateTransaction(c *fiber.Ctx) error {
	var req CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tx := models.Transaction{
		FromAccountID:  req.FromAccountID,
		ToAccountID:    req.ToAccountID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		IdempotencyKey: req.IdempotencyKey,
	}

	created, err := h.domain.CreateUsecase.Execute(c.UserContext(), tx)
	if err != nil {
		h.logger.Errorw("failed to create transaction", "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(201).JSON(mapToResponse(created))
}

func (h *handler) GetTransaction(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid transaction ID"})
	}

	tx, err := h.domain.GetUsecase.Execute(c.Context(), id)
	if err != nil {
		h.logger.Errorw("failed to get transaction", "id", id, "error", err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}
	if tx == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Transaction not found"})
	}

	return c.Status(200).JSON(mapToResponse(tx))
}

func mapToResponse(tx *models.Transaction) TransactionResponse {
	resp := TransactionResponse{
		ID:             tx.ID,
		FromAccountID:  tx.FromAccountID,
		ToAccountID:    tx.ToAccountID,
		Amount:         tx.Amount,
		Currency:       tx.Currency,
		Status:         tx.Status,
		IdempotencyKey: tx.IdempotencyKey,
		CreatedAt:      tx.CreatedAt.String(),
	}
	return resp
}
