package service

import (
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/repo"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/account"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/ledger"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/outbox"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/users"
	"github.com/w0ikid/yarmaq/pkg/zitadel"
	"go.uber.org/zap"
)

type Service struct {
	AccountService account.Service
	LedgerService  ledger.Service
	OutboxService  outbox.Service
	UserService    users.Service
}

func New(repositories *repo.Repository, zitadelClient *zitadel.Client, logger *zap.SugaredLogger) *Service {
	logger = logger.Named("service")
	return &Service{
		AccountService: account.NewService(repositories.Account, repositories.Ledger, repositories.Outbox, logger),
		LedgerService:  ledger.NewService(repositories.Ledger, logger),
		OutboxService:  outbox.NewService(repositories.Outbox, logger),
		UserService:    users.NewService(zitadelClient, logger),
	}
}
