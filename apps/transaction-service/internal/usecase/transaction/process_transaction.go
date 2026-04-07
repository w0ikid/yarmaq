package transaction

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/account"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/saga"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/transaction"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type ProcessTransactionSagaUsecase struct {
	usecase.BaseUsecase
	SagaService interface {
		CreateStep(ctx context.Context, step models.SagaStep) (*models.SagaStep, error)
		UpdateStepStatus(ctx context.Context, stepID uuid.UUID, status string, errStr *string) error
	}
	TransactionService interface {
		UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	}
	AccountService interface {
		Hold(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error
		Deposit(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error
		Refund(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error
	}
}

func NewProcessTransactionSagaUsecase(
	base usecase.BaseUsecase,
	sagaService saga.Service,
	transactionService transaction.Service,
	accountService account.Service,
) ProcessTransactionSagaUsecase {
	return ProcessTransactionSagaUsecase{
		BaseUsecase:        base,
		SagaService:        sagaService,
		TransactionService: transactionService,
		AccountService:     accountService,
	}
}

func (uc *ProcessTransactionSagaUsecase) Execute(ctx context.Context, event models.TransactionCreatedEvent) error {
	txID := uuid.MustParse(event.ID)

	// step 1: Hold
	holdStep, err := uc.SagaService.CreateStep(ctx, models.SagaStep{
		TransactionID: txID,
		StepName:      models.SagaStepNameHold,
		Status:        models.SagaStatusPending,
	})
	if err != nil {
		return err
	}
	
	err = uc.AccountService.Hold(ctx, event.FromAccountID, txID, event.Amount)
	if err != nil {
		errStr := err.Error()
		uc.SagaService.UpdateStepStatus(ctx, holdStep.ID, models.SagaStatusFailed, &errStr)
		uc.TransactionService.UpdateStatus(ctx, txID, models.TransactionStatusFailed)
		return err
	}
	uc.SagaService.UpdateStepStatus(ctx, holdStep.ID, models.SagaStatusCompleted, nil)
	uc.TransactionService.UpdateStatus(ctx, txID, models.TransactionStatusHolding)

	// step 2: Deposit
	depositStep, err := uc.SagaService.CreateStep(ctx, models.SagaStep{
		TransactionID: txID,
		StepName:      models.SagaStepNameDeposit,
		Status:        models.SagaStatusPending,
	})
	if err != nil {
		// refund hold
		uc.compensateHold(ctx, event.FromAccountID, txID, event.Amount)
		uc.TransactionService.UpdateStatus(ctx, txID, models.TransactionStatusFailed)
		return err
	}

	uc.TransactionService.UpdateStatus(ctx, txID, models.TransactionStatusDepositing)

	err = uc.AccountService.Deposit(ctx, event.ToAccountID, txID, event.TargetAmount)
	if err != nil {
		errStr := err.Error()
		if err := uc.SagaService.UpdateStepStatus(ctx, depositStep.ID, models.SagaStatusFailed, &errStr); err != nil {
			uc.Logger.Warnw("failed to update deposit step status", "error", err)
		}
		uc.compensateHold(ctx, event.FromAccountID, txID, event.Amount)
		if err := uc.TransactionService.UpdateStatus(ctx, txID, models.TransactionStatusFailed); err != nil {
			uc.Logger.Warnw("failed to update transaction status to failed", "error", err)
		}
		return err
	}

	if err := uc.SagaService.UpdateStepStatus(ctx, depositStep.ID, models.SagaStatusCompleted, nil); err != nil {
		uc.Logger.Warnw("failed to update deposit step status", "error", err)
	}
	if err := uc.TransactionService.UpdateStatus(ctx, txID, models.TransactionStatusCompleted); err != nil {
		uc.Logger.Warnw("failed to update transaction status to completed", "error", err)
	}

	return nil
}

func (uc *ProcessTransactionSagaUsecase) compensateHold(ctx context.Context, accountID string, txID uuid.UUID, amount int64) {
	refundStep, err := uc.SagaService.CreateStep(ctx, models.SagaStep{
		TransactionID: txID,
		StepName:      models.SagaStepNameRefund,
		Status:        models.SagaStatusPending,
	})
	if err != nil {
		uc.Logger.Errorw("failed to create refund step", "error", err)
		return
	}

	err = uc.AccountService.Refund(ctx, accountID, txID, amount)
	if err != nil {
		errStr := err.Error()
		uc.Logger.Errorw("compensation failed", "error", err)
		uc.SagaService.UpdateStepStatus(ctx, refundStep.ID, models.SagaStatusFailed, &errStr)
		return
	}

	uc.SagaService.UpdateStepStatus(ctx, refundStep.ID, models.SagaStatusCompleted, nil)
}
