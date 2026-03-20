package consumers

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type TransactionCreatedHandler struct {
	sagaUsecase interface {
		Execute(ctx context.Context, event models.TransactionCreatedEvent) error
	}
	logger *zap.SugaredLogger
}

func NewTransactionCreatedHandler(sagaUsecase interface {
	Execute(ctx context.Context, event models.TransactionCreatedEvent) error
}, logger *zap.SugaredLogger) *TransactionCreatedHandler {
	return &TransactionCreatedHandler{
		sagaUsecase: sagaUsecase,
		logger:      logger,
	}
}

func (h *TransactionCreatedHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var event models.TransactionCreatedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.logger.Errorw("failed to unmarshal transaction.created", "error", err)
		return err
	}

	h.logger.Infow("transaction.created received", "id", event.ID)

	if err := h.sagaUsecase.Execute(ctx, event); err != nil {
		h.logger.Errorw("saga failed", "id", event.ID, "error", err)
		return err
	}

	return nil
}
