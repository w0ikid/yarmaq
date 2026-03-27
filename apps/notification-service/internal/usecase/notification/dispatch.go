package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service/notification"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type DispatchNotificationUsecase struct {
	usecase.BaseUsecase
	NotificationService interface {
		Create(ctx context.Context, notification models.Notification) (*models.Notification, error)
		Send(ctx context.Context, notificationID uuid.UUID) (*models.Notification, error)
	}
}

func NewDispatchNotificationUsecase(base usecase.BaseUsecase, notificationService notification.Service) DispatchNotificationUsecase {
	return DispatchNotificationUsecase{
		BaseUsecase:         base,
		NotificationService: notificationService,
	}
}

func (uc *DispatchNotificationUsecase) Execute(ctx context.Context, notification models.Notification) (*models.Notification, error) {
	uc.Logger.Infow("starting DispatchNotificationUsecase execution", "user_id", notification.UserID, "type", notification.Type)

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
		uc.Logger.Errorw("failed to dispatch notification", "id", created.ID, "error", err)
		return nil, err
	}

	uc.Logger.Infow("DispatchNotificationUsecase executed successfully", "id", sent.ID, "status", sent.Status)
	return sent, nil
}
