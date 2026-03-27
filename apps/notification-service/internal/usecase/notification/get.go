package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service/notification"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type GetNotificationUsecase struct {
	usecase.BaseUsecase
	NotificationService interface {
		GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error)
	}
}

func NewGetNotificationUsecase(base usecase.BaseUsecase, notificationService notification.Service) GetNotificationUsecase {
	return GetNotificationUsecase{
		BaseUsecase:         base,
		NotificationService: notificationService,
	}
}

func (uc *GetNotificationUsecase) Execute(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	uc.Logger.Infow("fetching notification", "id", id)
	return uc.NotificationService.GetByID(ctx, id)
}
