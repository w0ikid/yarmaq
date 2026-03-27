package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service/notification"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type SendNotificationUsecase struct {
	usecase.BaseUsecase
	NotificationService interface {
		Create(ctx context.Context, notification models.Notification) (*models.Notification, error)
		Send(ctx context.Context, notificationID uuid.UUID) (*models.Notification, error)
	}
}

func NewSendNotificationUsecase(base usecase.BaseUsecase, notificationService notification.Service) SendNotificationUsecase {
	return SendNotificationUsecase{
		BaseUsecase:         base,
		NotificationService: notificationService,
	}
}

func (uc *SendNotificationUsecase) Execute(ctx context.Context, notification models.Notification) (*models.Notification, error) {
	uc.Logger.Infow("starting SendNotificationUsecase execution", "user_id", notification.UserID, "type", notification.Type)

	txCtx, err := uc.Tx.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer uc.Tx.FinalizeTransaction(txCtx, &err)

	created, err := uc.NotificationService.Create(txCtx, notification)
	if err != nil {
		uc.Logger.Errorw("failed to create notification before send", "user_id", notification.UserID, "type", notification.Type, "error", err)
		return nil, err
	}

	sent, err := uc.NotificationService.Send(txCtx, created.ID)
	if err != nil {
		uc.Logger.Errorw("failed to send notification", "id", created.ID, "error", err)
		return nil, err
	}

	uc.Logger.Infow("SendNotificationUsecase executed successfully", "id", sent.ID, "status", sent.Status)
	return sent, nil
}
