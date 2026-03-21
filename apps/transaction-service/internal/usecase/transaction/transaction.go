package transaction

import (
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/account"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/outbox"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/saga"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/transaction"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/usecase"
)

type TransactionDomain struct {
	CreateUsecase CreateTransactionUsecase
	GetUsecase    GetTransactionUsecase

	ProcessSagaUsecase ProcessTransactionSagaUsecase
}

func NewDomain(baseusecase usecase.BaseUsecase, transactionService transaction.Service, outboxService outbox.Service, sagaService saga.Service, accountService account.Service) TransactionDomain {
	baseusecase.Logger = baseusecase.Logger.Named("transaction_domain")
	return TransactionDomain{
		CreateUsecase:      NewCreateTransactionUsecase(baseusecase, transactionService, outboxService),
		GetUsecase:         NewGetTransactionUsecase(baseusecase, transactionService),
		ProcessSagaUsecase: NewProcessTransactionSagaUsecase(baseusecase, sagaService, transactionService, accountService),
	}
}
