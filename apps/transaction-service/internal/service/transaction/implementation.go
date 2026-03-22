package transaction

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/ctxkeys"
	"github.com/w0ikid/yarmaq/pkg/errs"
	"github.com/w0ikid/yarmaq/pkg/httpclient/accounts"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, transaction models.Transaction) (*models.Transaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type implementation struct {
	repo           TransactionRepo
	accountsClient *accounts.Client
	logger         *zap.SugaredLogger
}

func NewService(repo TransactionRepo, accountsClient *accounts.Client, logger *zap.SugaredLogger) Service {
	return &implementation{
		repo:           repo,
		accountsClient: accountsClient,
		logger:         logger.Named("transaction_service"),
	}
}

func (s *implementation) Create(ctx context.Context, transaction models.Transaction) (*models.Transaction, error) {
	s.logger.Infow("creating transaction", "to_account_number", transaction.ToAccountNumber, "amount", transaction.Amount, "currency", transaction.Currency)

	if transaction.Amount <= 0 {
		return nil, fmt.Errorf("%w: transaction amount must be positive", errs.ErrValidation)
	}

	if transaction.ToAccountNumber == "" {
		return nil, fmt.Errorf("%w: to_account_number is required", errs.ErrValidation)
	}

	userID := ctxkeys.GetUserID(ctx)
	if userID == "" {
		return nil, fmt.Errorf("%w: user is not authenticated", errs.ErrUnauthorized)
	}

	fromAccount, err := s.accountsClient.GetAccountByUserIDAndCurrency(ctx, userID, transaction.Currency)
	if err != nil {
		return nil, fmt.Errorf("get from_account: %w", err)
	}
	s.logger.Infow("from_account fetched", "user_id", userID, "currency", transaction.Currency)

	if fromAccount == nil {
		return nil, fmt.Errorf("%w: from_account not found for user %s and currency %s", errs.ErrNotFound, userID, transaction.Currency)
	}
	transaction.FromAccountID = fromAccount.ID

	s.logger.Infow("fetching to_account", "number", transaction.ToAccountNumber, "currency", transaction.Currency)
	toAccount, err := s.accountsClient.GetAccountByNumberAndCurrency(ctx, transaction.ToAccountNumber, transaction.Currency)
	if err != nil {
		return nil, fmt.Errorf("get to_account: %w", err)
	}
	s.logger.Infow("to_account fetched", "number", transaction.ToAccountNumber, "currency", transaction.Currency)

	if toAccount == nil {
		return nil, fmt.Errorf("%w: to_account not found for number %s and currency %s", errs.ErrNotFound, transaction.ToAccountNumber, transaction.Currency)
	}
	transaction.ToAccountID = toAccount.ID

	if transaction.FromAccountID == transaction.ToAccountID {
		return nil, fmt.Errorf("%w: from and to accounts cannot be the same", errs.ErrValidation)
	}

	if fromAccount.Currency != transaction.Currency || toAccount.Currency != transaction.Currency {
		return nil, fmt.Errorf("%w: accounts have different currencies: %s vs %s", errs.ErrValidation, fromAccount.Currency, toAccount.Currency)
	}

	transaction.Currency = fromAccount.Currency
	transaction.Status = models.TransactionStatusPending

	if transaction.IdempotencyKey != "" {
		existing, err := s.repo.GetByIdempotencyKey(ctx, transaction.IdempotencyKey)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			s.logger.Infow("transaction already exists (idempotency)", "id", existing.ID)
			return existing, nil
		}
	}

	return s.repo.Create(ctx, transaction)
}

func (s *implementation) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *implementation) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	s.logger.Infow("updating transaction status", "id", id, "status", status)
	tx, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if tx == nil {
		return fmt.Errorf("%w: transaction not found: %s", errs.ErrNotFound, id)
	}

	tx.Status = status
	_, err = s.repo.Update(ctx, *tx)
	return err
}
