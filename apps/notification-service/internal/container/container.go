package container

import (
	"context"

	"github.com/w0ikid/yarmaq/apps/notification-service/internal/repo"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/usecase"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/usecase/notification"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/usecase/outbox"
	"github.com/w0ikid/yarmaq/pkg/httpclient/accounts"
	"github.com/w0ikid/yarmaq/pkg/smtpclient"
	"github.com/w0ikid/yarmaq/pkg/zitadel"
	"go.uber.org/zap"
)

type Container struct {
	logger *zap.SugaredLogger

	Services *service.Service

	NotificationDomain notification.NotificationDomain
	OutboxDomain       outbox.OutboxDomain
}

func NewContainer(
	ctx context.Context,
	repositories *repo.Repository,
	zitadelClient *zitadel.Client,
	accountsClient *accounts.Client,
	smtpClient *smtpclient.Client,
	logger *zap.SugaredLogger,

) *Container {
	logger = logger.Named("container")

	services := service.New(repositories, zitadelClient, accountsClient, smtpClient, logger)

	baseusecase := usecase.BaseUsecase{
		Logger: logger.Named("base_usecase"),
		Tx:     repositories.ContextTransaction,
	}

	return &Container{
		logger:             logger,
		Services:           services,
		NotificationDomain: notification.NewDomain(baseusecase, services.NotificationService),
		OutboxDomain:       outbox.NewDomain(baseusecase, services.OutboxService),
	}
}
