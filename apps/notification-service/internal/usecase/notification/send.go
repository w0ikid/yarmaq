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
		Send(ctx context.Context, notificationID uuid.UUID) (*models.Notification, error)
	}
}

func NewSendNotificationUsecase(base usecase.BaseUsecase, notificationService notification.Service) SendNotificationUsecase {
	return SendNotificationUsecase{
		BaseUsecase:         base,
		NotificationService: notificationService,
	}
}

func (uc *SendNotificationUsecase) Execute(ctx context.Context, notificationID uuid.UUID) (*models.Notification, error) {
	uc.Logger.Infow("starting SendNotificationUsecase execution", "id", notificationID)
	sent, err := uc.NotificationService.Send(ctx, notificationID)
	if err != nil {
		uc.Logger.Errorw("failed to send notification", "id", notificationID, "error", err)
		return nil, err
	}

	uc.Logger.Infow("SendNotificationUsecase executed successfully", "id", sent.ID, "status", sent.Status)
	return sent, nil
}
