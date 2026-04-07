package transaction

import (
	"context"
	"encoding/json"

	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/outbox"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/service/transaction"
	"github.com/w0ikid/yarmaq/apps/transaction-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type CreateTransactionUsecase struct {
	usecase.BaseUsecase
	TransactionService interface {
		Create(ctx context.Context, transaction models.Transaction) (*models.Transaction, error)
	}
	OutboxService interface {
		Create(ctx context.Context, event models.Outbox) (*models.Outbox, error)
	}
}

func NewCreateTransactionUsecase(base usecase.BaseUsecase, transactionService transaction.Service, outboxService outbox.Service) CreateTransactionUsecase {
	return CreateTransactionUsecase{
		BaseUsecase:        base,
		TransactionService: transactionService,
		OutboxService:      outboxService,
	}
}

func (uc *CreateTransactionUsecase) Execute(ctx context.Context, transaction models.Transaction) (*models.Transaction, error) {
	uc.Logger.Infow("starting CreateTransactionUsecase execution", "idempotency_key", transaction.IdempotencyKey)

	txCtx, err := uc.Tx.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer uc.Tx.FinalizeTransaction(txCtx, &err)

	created, err := uc.TransactionService.Create(txCtx, transaction)
	if err != nil {
		uc.Logger.Errorw("failed to create transaction", "error", err)
		return nil, err
	}

	payload, err := json.Marshal(models.TransactionCreatedEvent{
		ID:            created.ID.String(),
		Type:          created.Type,
		FromAccountID: created.FromAccountID.String(),
		ToAccountID:   created.ToAccountID.String(),
		Amount:        created.Amount,
		Currency:      created.Currency,
		TargetAmount:   *created.TargetAmount,
		TargetCurrency: *created.TargetCurrency,
		ExchangeRate:   *created.ExchangeRate,
	})
	if err != nil {
		uc.Logger.Errorw("failed to marshal transaction created event", "error", err)
		return nil, err
	}

	_, err = uc.OutboxService.Create(txCtx, models.Outbox{
		EventType:   "transaction.created",
		Payload:     payload,
		AggregateID: created.ID,
	})
	if err != nil {
		uc.Logger.Errorw("failed to create outbox event", "error", err)
		return nil, err
	}

	uc.Logger.Infow("CreateTransactionUsecase executed successfully", "id", created.ID)
	return created, nil
}
