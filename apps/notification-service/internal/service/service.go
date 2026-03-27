package service

import (
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/repo"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service/account"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service/notification"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service/outbox"
	"github.com/w0ikid/yarmaq/apps/notification-service/internal/service/transaction"
	"github.com/w0ikid/yarmaq/pkg/httpclient/accounts"
	"github.com/w0ikid/yarmaq/pkg/smtpclient"
	"github.com/w0ikid/yarmaq/pkg/zitadel"
	"go.uber.org/zap"
)

type Service struct {
	AccountService      account.Service
	NotificationService notification.Service
	OutboxService       outbox.Service
	TransactionService  transaction.Service
}

func New(repositories *repo.Repository, zitadelClient *zitadel.Client, accountsClient *accounts.Client, smtpClient *smtpclient.Client, logger *zap.SugaredLogger) *Service {
	logger = logger.Named("service")
	return &Service{
		AccountService:      account.NewService(accountsClient, logger),
		NotificationService: notification.NewService(repositories.Notification, smtpClient, logger),
		OutboxService:       outbox.NewService(repositories.Outbox, logger),
	}
}
