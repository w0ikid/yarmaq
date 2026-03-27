package notification

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service/notification"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type GetNotificationsByUserUsecase struct {
	usecase.BaseUsecase
	NotificationService interface {
		GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.Notification, error)
	}
}

func NewGetNotificationsByUserUsecase(base usecase.BaseUsecase, notificationService notification.Service) GetNotificationsByUserUsecase {
	return GetNotificationsByUserUsecase{
		BaseUsecase:         base,
		NotificationService: notificationService,
	}
}

func (uc *GetNotificationsByUserUsecase) Execute(ctx context.Context, userID uuid.UUID) ([]models.Notification, error) {
	uc.Logger.Infow("fetching notifications by user", "user_id", userID)
	return uc.NotificationService.GetByUserID(ctx, userID)
}
