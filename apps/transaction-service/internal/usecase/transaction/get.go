package transaction

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/transaction"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type GetTransactionUsecase struct {
	usecase.BaseUsecase
	TransactionService interface {
		GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	}
}

func NewGetTransactionUsecase(base usecase.BaseUsecase, transactionService transaction.Service) GetTransactionUsecase {
	return GetTransactionUsecase{
		BaseUsecase:        base,
		TransactionService: transactionService,
	}
}

func (uc *GetTransactionUsecase) Execute(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	uc.Logger.Infow("fetching transaction", "id", id)
	return uc.TransactionService.GetByID(ctx, id)
}
