package transaction

import (
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/transaction"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/usecase"
)

type TransactionDomain struct {
	CreateUsecase      CreateTransactionUsecase
	GetUsecase         GetTransactionUsecase
}

func NewDomain(baseusecase usecase.BaseUsecase, transactionService transaction.Service) TransactionDomain {
	baseusecase.Logger = baseusecase.Logger.Named("transaction_domain")
	return TransactionDomain{
		CreateUsecase: CreateTransactionUsecase{
			BaseUsecase:        baseusecase,
			TransactionService: transactionService,
		},
		GetUsecase: GetTransactionUsecase{
			BaseUsecase:        baseusecase,
			TransactionService: transactionService,
		},
	}
}
