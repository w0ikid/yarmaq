package service

import (
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/repo"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/account"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/outbox"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/saga"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/transaction"
	"github.com/w0ikid/yarmaq/pkg/exchange"
	"github.com/w0ikid/yarmaq/pkg/httpclient/accounts"
	"github.com/w0ikid/yarmaq/pkg/zitadel"
	"go.uber.org/zap"
)

type Service struct {
	AccountService     account.Service
	OutboxService      outbox.Service
	TransactionService transaction.Service
	SagaService        saga.Service
}

func New(repositories *repo.Repository, zitadelClient *zitadel.Client, accountsClient *accounts.Client, logger *zap.SugaredLogger) *Service {
	logger = logger.Named("service")
	exchangeSvc := exchange.NewService()
	return &Service{
		AccountService:     account.NewService(accountsClient, logger),
		OutboxService:      outbox.NewService(repositories.Outbox, logger),
		TransactionService: transaction.NewService(repositories.Transaction, accountsClient, exchangeSvc, logger),
		SagaService:        saga.NewService(repositories.SagaStep, logger),
	}
}
