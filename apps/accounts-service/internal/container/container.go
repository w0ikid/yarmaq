package container

import (
	"context"

	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/repo"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/account"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/ledger"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/outbox"
	"github.com/w0ikid/yarmaq/pkg/zitadel"
	"go.uber.org/zap"
)

type Container struct {
	logger *zap.SugaredLogger

	Services *service.Service

	AccountDomain account.AccountDomain
	LedgerDomain  ledger.LedgerDomain
	OutboxDomain  outbox.OutboxDomain
}

func NewContainer(
	ctx context.Context,
	repositories *repo.Repository,
	zitadelClient *zitadel.Client,
	logger *zap.SugaredLogger,

) *Container {
	logger = logger.Named("container")

	services := service.New(repositories, zitadelClient, logger)

	baseusecase := usecase.BaseUsecase{
		Logger: logger.Named("base_usecase"),
		Tx:     repositories.ContextTransaction,
	}

	return &Container{
		logger: logger,

		Services: services,

		AccountDomain: account.NewDomain(baseusecase, services.AccountService, services.OutboxService, services.UserService),
		LedgerDomain:  ledger.NewDomain(baseusecase, services.LedgerService),
		OutboxDomain:  outbox.NewDomain(baseusecase, services.OutboxService),
	}

}
