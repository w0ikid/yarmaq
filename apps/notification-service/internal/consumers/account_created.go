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

type AccountCreatedHandler struct {
	dispatchNotificationUsecase interface {
		Execute(ctx context.Context, notification models.Notification) (*models.Notification, error)
	}
	emailResolver UserEmailResolver
	logger        *zap.SugaredLogger
}

func NewAccountCreatedHandler(dispatchNotificationUsecase interface {
	Execute(ctx context.Context, notification models.Notification) (*models.Notification, error)
}, emailResolver UserEmailResolver, logger *zap.SugaredLogger) *AccountCreatedHandler {
	return &AccountCreatedHandler{
		dispatchNotificationUsecase: dispatchNotificationUsecase,
		emailResolver:               emailResolver,
		logger:                      logger,
	}
}

func (h *AccountCreatedHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var event models.AccountCreatedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.logger.Errorw("failed to unmarshal account.created", "error", err)
		return err
	}

	h.logger.Infow("account.created received", "account_id", event.ID, "user_id", event.UserID)

	userID, err := uuid.Parse(event.UserID)
	if err != nil {
		return fmt.Errorf("parse account.created user_id: %w", err)
	}

	email, err := h.emailResolver.ResolveEmail(ctx, event.UserID)
	if err != nil {
		h.logger.Errorw("failed to resolve recipient email", "user_id", event.UserID, "error", err)
		return err
	}

	notification, err := h.dispatchNotificationUsecase.Execute(ctx, models.Notification{
		UserID:   userID,
		Type:     models.TypeAccountCreated,
		Channel:  models.ChannelEmail,
		Subject:  "Account created",
		Body:     fmt.Sprintf("<p>Your account %s was created successfully.</p>", event.ID),
		Metadata: map[string]any{"email": email, "account_id": event.ID},
	})
	if err != nil {
		h.logger.Errorw("failed to dispatch account created notification", "account_id", event.ID, "user_id", event.UserID, "error", err)
		return err
	}

	h.logger.Infow("account created notification sent", "notification_id", notification.ID, "user_id", event.UserID, "status", notification.Status)
	return nil
}
