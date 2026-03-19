package transaction

import (
	"context"

	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type CreateTransactionUsecase struct {
	usecase.BaseUsecase
	TransactionService interface {
		Create(ctx context.Context, transaction models.Transaction) (*models.Transaction, error)
	}
}

func (uc *CreateTransactionUsecase) Execute(ctx context.Context, transaction models.Transaction) (*models.Transaction, error) {
	uc.Logger.Infow("starting CreateTransactionUsecase execution", "idempotency_key", transaction.IdempotencyKey)

	txCtx, err := uc.Tx.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer uc.Tx.FinalizeTransaction(txCtx, &err)

	transaction.Status = models.TransactionStatusPending
	created, err := uc.TransactionService.Create(txCtx, transaction)
	if err != nil {
		uc.Logger.Errorw("failed to create transaction", "error", err)
		return nil, err
	}

	uc.Logger.Infow("CreateTransactionUsecase executed successfully", "id", created.ID)
	return created, nil
}
