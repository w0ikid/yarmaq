package transaction

import (
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

func NewDomain(baseusecase usecase.BaseUsecase, transactionService transaction.Service, outboxService outbox.Service, sagaService saga.Service) TransactionDomain {
	baseusecase.Logger = baseusecase.Logger.Named("transaction_domain")
	return TransactionDomain{
		CreateUsecase: CreateTransactionUsecase{
			BaseUsecase:        baseusecase,
			TransactionService: transactionService,
			OutboxService:      outboxService,
		},
		GetUsecase: GetTransactionUsecase{
			BaseUsecase:        baseusecase,
			TransactionService: transactionService,
		},
		ProcessSagaUsecase: ProcessTransactionSagaUsecase{
			BaseUsecase:        baseusecase,
			SagaService:        sagaService,
			TransactionService: transactionService,
		},
	}
}
