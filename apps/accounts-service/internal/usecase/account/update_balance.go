package account

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/account"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/outbox"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type UpdateBalanceUsecase struct {
	usecase.BaseUsecase
	AccountService interface {
		UpdateBalance(ctx context.Context, accountID uuid.UUID, amount int64, operationType string, referenceID *uuid.UUID) error
	}
	OutboxService interface {
		Create(ctx context.Context, event models.Outbox) (*models.Outbox, error)
	}
}

func NewUpdateBalanceUsecase(base usecase.BaseUsecase, accountService account.Service, outboxService outbox.Service) UpdateBalanceUsecase {
	return UpdateBalanceUsecase{
		BaseUsecase:    base,
		AccountService: accountService,
		OutboxService:  outboxService,
	}
}

func (uc *UpdateBalanceUsecase) Execute(ctx context.Context, accountID uuid.UUID, amount int64, operationType string, referenceID *uuid.UUID) error {
	uc.Logger.Infow("starting UpdateBalanceUsecase", "account_id", accountID, "amount", amount)

	txCtx, err := uc.Tx.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer uc.Tx.FinalizeTransaction(txCtx, &err)

	if err = uc.AccountService.UpdateBalance(txCtx, accountID, amount, operationType, referenceID); err != nil {
		uc.Logger.Errorw("failed to update balance", "account_id", accountID, "error", err)
		return err
	}

	payload, _ := json.Marshal(models.BalanceUpdatedPayload{
		AccountID:     accountID,
		Amount:        amount,
		OperationType: operationType,
		ReferenceID:   referenceID,
	})

	if _, err = uc.OutboxService.Create(txCtx, models.Outbox{
		EventType:   "account.balance.updated",
		AggregateID: accountID,
		Payload:     payload,
	}); err != nil {
		uc.Logger.Errorw("failed to create outbox event", "account_id", accountID, "error", err)
		return err
	}

	uc.Logger.Infow("UpdateBalanceUsecase executed successfully", "account_id", accountID)
	return nil
}
