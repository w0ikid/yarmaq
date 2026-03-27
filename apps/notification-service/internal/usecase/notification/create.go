package notification

import (
	"context"

	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service/notification"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type CreateNotificationUsecase struct {
	usecase.BaseUsecase
	NotificationService interface {
		Create(ctx context.Context, notification models.Notification) (*models.Notification, error)
	}
}

func NewCreateNotificationUsecase(base usecase.BaseUsecase, notificationService notification.Service) CreateNotificationUsecase {
	return CreateNotificationUsecase{
		BaseUsecase:         base,
		NotificationService: notificationService,
	}
}

func (uc *CreateNotificationUsecase) Execute(ctx context.Context, notification models.Notification) (*models.Notification, error) {
	uc.Logger.Infow("starting CreateNotificationUsecase execution", "user_id", notification.UserID, "type", notification.Type)

	txCtx, err := uc.Tx.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer uc.Tx.FinalizeTransaction(txCtx, &err)

	created, err := uc.NotificationService.Create(txCtx, notification)
	if err != nil {
		uc.Logger.Errorw("failed to create notification", "user_id", notification.UserID, "type", notification.Type, "error", err)
		return nil, err
	}

	uc.Logger.Infow("CreateNotificationUsecase executed successfully", "id", created.ID)
	return created, nil
}
