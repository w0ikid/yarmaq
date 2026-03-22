package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/w0ikid/yarmaq/pkg/ctxkeys"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type AccountCreatedHandler struct {
	createTransactionUsecase interface {
		Execute(ctx context.Context, transaction models.Transaction) (*models.Transaction, error)
	}
	logger *zap.SugaredLogger
}

func NewAccountCreatedHandler(createTransactionUsecase interface {
	Execute(ctx context.Context, transaction models.Transaction) (*models.Transaction, error)
}, logger *zap.SugaredLogger) *AccountCreatedHandler {
	return &AccountCreatedHandler{
		createTransactionUsecase: createTransactionUsecase,
		logger:                   logger,
	}
}

func (h *AccountCreatedHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var event models.AccountCreatedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.logger.Errorw("failed to unmarshal account.created", "error", err)
		return err
	}

	h.logger.Infow("account.created received", "id", event.ID, "user_id", event.UserID)

	txCtx := ctxkeys.WithUserContext(ctx, event.UserID, nil)
	created, err := h.createTransactionUsecase.Execute(txCtx, models.Transaction{
		Type:           models.TransactionTypeDeposit,
		Amount:         1000,
		Currency:       "KZT",
		IdempotencyKey: fmt.Sprintf("account-created:%s:kzt-bonus", event.ID),
	})
	if err != nil {
		h.logger.Errorw("failed to create welcome deposit transaction", "account_id", event.ID, "user_id", event.UserID, "error", err)
		return err
	}

	h.logger.Infow("welcome deposit transaction created", "transaction_id", created.ID, "user_id", event.UserID, "amount", created.Amount, "currency", created.Currency)
	return nil
}
