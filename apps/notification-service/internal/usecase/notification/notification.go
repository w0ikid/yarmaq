package notification

import (
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service/notification"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/usecase"
)

type NotificationDomain struct {
	CreateUsecase    CreateNotificationUsecase
	GetUsecase       GetNotificationUsecase
	GetByUserUsecase GetNotificationsByUserUsecase
	SendUsecase      SendNotificationUsecase
	DispatchUsecase  DispatchNotificationUsecase
}

func NewDomain(baseusecase usecase.BaseUsecase, notificationService notification.Service) NotificationDomain {
	baseusecase.Logger = baseusecase.Logger.Named("notification_domain")

	return NotificationDomain{
		CreateUsecase:    NewCreateNotificationUsecase(baseusecase, notificationService),
		GetUsecase:       NewGetNotificationUsecase(baseusecase, notificationService),
		GetByUserUsecase: NewGetNotificationsByUserUsecase(baseusecase, notificationService),
		SendUsecase:      NewSendNotificationUsecase(baseusecase, notificationService),
		DispatchUsecase:  NewDispatchNotificationUsecase(baseusecase, notificationService),
	}
}
