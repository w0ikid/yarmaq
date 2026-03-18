package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
)

type UpdateBalanceUsecase struct {
	usecase.BaseUsecase
	AccountService interface {
		UpdateBalance(ctx context.Context, accountID uuid.UUID, amount int64, operationType string, referenceID *uuid.UUID) error
	}
}

func (uc *UpdateBalanceUsecase) Execute(ctx context.Context, accountID uuid.UUID, amount int64, operationType string, referenceID *uuid.UUID) error {
	uc.Logger.Infow("starting UpdateBalanceUsecase execution", "account_id", accountID, "amount", amount)

	txCtx, err := uc.Tx.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer uc.Tx.FinalizeTransaction(txCtx, &err)

	err = uc.AccountService.UpdateBalance(txCtx, accountID, amount, operationType, referenceID)
	if err != nil {
		uc.Logger.Errorw("failed to update balance", "account_id", accountID, "error", err)
		return err
	}

	uc.Logger.Infow("UpdateBalanceUsecase executed successfully", "account_id", accountID)
	return nil
}
