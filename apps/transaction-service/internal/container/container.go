package container

import (
	"context"

	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/repo"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/usecase"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/usecase/outbox"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/usecase/transaction"
	"github.com/w0ikid/yarmaq/pkg/httpclient/accounts"
	"github.com/w0ikid/yarmaq/pkg/zitadel"
	"go.uber.org/zap"
)

type Container struct {
	logger *zap.SugaredLogger

	Services *service.Service

	OutboxDomain      outbox.OutboxDomain
	TransactionDomain transaction.TransactionDomain
}

func NewContainer(
	ctx context.Context,
	repositories *repo.Repository,
	zitadelClient *zitadel.Client,
	accountsClient *accounts.Client,
	logger *zap.SugaredLogger,

) *Container {
	logger = logger.Named("container")

	services := service.New(repositories, zitadelClient, accountsClient, logger)

	baseusecase := usecase.BaseUsecase{
		Logger: logger.Named("base_usecase"),
		Tx:     repositories.ContextTransaction,
	}

	return &Container{
		logger: logger,
		Services: services,
		OutboxDomain:      outbox.NewDomain(baseusecase, services.OutboxService),
		TransactionDomain: transaction.NewDomain(baseusecase, services.TransactionService, services.OutboxService, services.SagaService, services.AccountService),
	}
}
