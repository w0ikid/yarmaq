package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type TransactionCreatedHandler struct {
	dispatchNotificationUsecase interface {
		Execute(ctx context.Context, notification models.Notification) (*models.Notification, error)
	}
	accountResolver AccountUserResolver
	emailResolver   UserEmailResolver
	logger          *zap.SugaredLogger
}

func NewTransactionCreatedHandler(dispatchNotificationUsecase interface {
	Execute(ctx context.Context, notification models.Notification) (*models.Notification, error)
}, accountResolver AccountUserResolver, emailResolver UserEmailResolver, logger *zap.SugaredLogger) *TransactionCreatedHandler {
	return &TransactionCreatedHandler{
		dispatchNotificationUsecase: dispatchNotificationUsecase,
		accountResolver:             accountResolver,
		emailResolver:               emailResolver,
		logger:                      logger,
	}
}

func (h *TransactionCreatedHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var event models.TransactionCreatedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.logger.Errorw("failed to unmarshal transaction.created", "error", err)
		return err
	}

	h.logger.Infow("transaction.created received", "id", event.ID)

	recipients := make(map[string]struct{})
	for _, accountID := range []string{event.FromAccountID, event.ToAccountID} {
		userID, err := h.accountResolver.ResolveUserID(ctx, accountID)
		if err != nil {
			h.logger.Errorw("failed to resolve account owner", "account_id", accountID, "transaction_id", event.ID, "error", err)
			return err
		}
		if userID == "" {
			continue
		}
		recipients[userID] = struct{}{}
	}

	for userIDStr := range recipients {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return fmt.Errorf("parse transaction user_id: %w", err)
		}

		email, err := h.emailResolver.ResolveEmail(ctx, userIDStr)
		if err != nil {
			h.logger.Errorw("failed to resolve recipient email", "user_id", userIDStr, "transaction_id", event.ID, "error", err)
			return err
		}

		notification, err := h.dispatchNotificationUsecase.Execute(ctx, models.Notification{
			UserID:   userID,
			Type:     models.TypeTransactionCreated,
			Channel:  models.ChannelEmail,
			Subject:  "Transaction created",
			Body:     fmt.Sprintf("<p>Transaction %s created: %s %d %s.</p>", event.ID, event.Type, event.Amount, event.Currency),
			Metadata: map[string]any{"email": email, "transaction_id": event.ID},
		})
		if err != nil {
			h.logger.Errorw("failed to dispatch transaction notification", "transaction_id", event.ID, "user_id", userIDStr, "error", err)
			return err
		}

		h.logger.Infow("transaction notification sent", "notification_id", notification.ID, "transaction_id", event.ID, "user_id", userIDStr, "status", notification.Status)
	}

	return nil
}
