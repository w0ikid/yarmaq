package account

import (
	"context"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
	"encoding/json"
)

type CreateAccountUsecase struct {
	usecase.BaseUsecase
	AccountService interface {
		Create(ctx context.Context, account models.Account) (*models.Account, error)
	}

	OutboxService interface {
		Create(ctx context.Context, event models.Outbox) (*models.Outbox, error)
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

	payload, err := json.Marshal(models.AccountCreatedEvent{
		ID:     createdAccount.ID.String(),
		UserID: createdAccount.UserID,
	})

	// outbox event
	_, err = uc.OutboxService.Create(txCtx, models.Outbox{
		EventType:   "account.created",
		Payload:     payload,
		AggregateID: createdAccount.ID,
	})
	if err != nil {
		uc.Logger.Errorw("failed to create outbox event", "user_id", account.UserID, "error", err)
		return nil, err
	}

	uc.Logger.Infow("CreateAccountUsecase executed successfully", "id", createdAccount.ID)
	return createdAccount, nil
}
