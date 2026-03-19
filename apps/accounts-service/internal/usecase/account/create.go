package account

import (
	"context"
	
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type CreateAccountUsecase struct {
	usecase.BaseUsecase
	AccountService interface {
		Create(ctx context.Context, account models.Account) (*models.Account, error)
	}
}

func (uc *CreateAccountUsecase) Execute(ctx context.Context, account models.Account) (*models.Account, error) {
	uc.Logger.Infow("starting CreateAccountUsecase execution", "user_id", account.UserID)

	txCtx, err := uc.Tx.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer uc.Tx.FinalizeTransaction(txCtx, &err)

	createdAccount, err := uc.AccountService.Create(txCtx, account)
	if err != nil {
		uc.Logger.Errorw("failed to create account", "user_id", account.UserID, "error", err)
		return nil, err
	}

	uc.Logger.Infow("CreateAccountUsecase executed successfully", "id", createdAccount.ID)
	return createdAccount, nil
}
