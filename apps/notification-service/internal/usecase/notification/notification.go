package notification

import (
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service/notification"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/usecase"
)

type NotificationDomain struct {
	SendUsecase SendNotificationUsecase
}

func NewDomain(baseusecase usecase.BaseUsecase, notificationService notification.Service) NotificationDomain {
	baseusecase.Logger = baseusecase.Logger.Named("notification_domain")

	return NotificationDomain{
		SendUsecase: NewSendNotificationUsecase(baseusecase, notificationService),
	}
}
