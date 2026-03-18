package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type Service interface {
	Create(ctx context.Context, account models.Account) (*models.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	GetByNumber(ctx context.Context, number string) (*models.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error)
	UpdateBalance(ctx context.Context, accountID uuid.UUID, amount int64, operationType string, referenceID *uuid.UUID) error
}

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
