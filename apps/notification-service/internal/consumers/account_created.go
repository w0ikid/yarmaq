package consumers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type AccountCreatedHandler struct {
	sendNotificationUsecase interface {
		Execute(ctx context.Context, notification models.Notification) (*models.Notification, error)
	}
	logger *zap.SugaredLogger
}

func NewAccountCreatedHandler(sendNotificationUsecase interface {
	Execute(ctx context.Context, notification models.Notification) (*models.Notification, error)
}, logger *zap.SugaredLogger) *AccountCreatedHandler {
	return &AccountCreatedHandler{
		sendNotificationUsecase: sendNotificationUsecase,
		logger:                  logger,
	}
}

func (h *AccountCreatedHandler) Handle(ctx context.Context, msg kafka.Message) error {
	var event models.AccountCreatedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		h.logger.Errorw("failed to unmarshal account.created", "error", err)
		return err
	}

	h.logger.Infow("account.created received", 
		"id", event.ID, 
		"user_id", event.UserID, 
		"number", event.Number, 
		"currency", event.Currency,
	)

	// 2. Map to Notification Model
	subject := fmt.Sprintf("New account created: %s", event.Number)
	body := fmt.Sprintf(`
		<div style="font-family: sans-serif; max-width: 600px; margin: 20px auto; padding: 20px; border: 1px solid #eee; border-radius: 10px; box-shadow: 0 4px 6px rgba(0,0,0,0.05);">
			<h2 style="color: #2d3748; margin-top: 0;">Account Created Successfully!</h2>
			<p style="color: #4a5568; line-height: 1.6;">Hello,</p>
			<p style="color: #4a5568; line-height: 1.6;">Your new account has been successfully opened and is ready for use.</p>
			<div style="background-color: #f7fafc; padding: 15px; border-radius: 8px; margin: 20px 0;">
				<table style="width: 100%%; border-collapse: collapse;">
					<tr>
						<td style="color: #718096; padding: 5px 0;">Account Number:</td>
						<td style="color: #2d3748; font-weight: bold; text-align: right;">%s</td>
					</tr>
					<tr>
						<td style="color: #718096; padding: 5px 0;">Currency:</td>
						<td style="color: #2d3748; font-weight: bold; text-align: right;">%s</td>
					</tr>
				</table>
			</div>
			<p style="color: #4a5568; line-height: 1.6;">Thank you for choosing Yarmaq.</p>
			<hr style="border: 0; border-top: 1px solid #edf2f7; margin: 20px 0;">
			<p style="color: #a0aec0; font-size: 12px; text-align: center;">This is an automated message, please do not reply.</p>
		</div>
	`, event.Number, event.Currency)

	notification, err := h.sendNotificationUsecase.Execute(ctx, models.Notification{
		UserID:   event.UserID,
		Type:     models.TypeAccountCreated,
		Channel:  models.ChannelEmail,
		Subject:  subject,
		Body:     body,
		Metadata: map[string]any{"email": event.Email, "account_id": event.ID},
	})
	if err != nil {
		h.logger.Errorw("failed to send account created notification", "account_id", event.ID, "user_id", event.UserID, "error", err)
		return err
	}

	h.logger.Infow("account created notification sent", "notification_id", notification.ID, "user_id", event.UserID, "status", notification.Status)
	return nil
}
