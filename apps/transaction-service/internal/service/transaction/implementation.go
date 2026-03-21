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
	repo    TransactionRepo
	accountsClient *accounts.Client
	logger  *zap.SugaredLogger
}

func NewService(repo TransactionRepo, accountsClient *accounts.Client, logger *zap.SugaredLogger) Service {
	return &implementation{
		repo:   repo,
		accountsClient: accountsClient,
		logger: logger.Named("transaction_service"),
	}
}

func (s *implementation) Create(ctx context.Context, transaction models.Transaction) (*models.Transaction, error) {
	s.logger.Infow("creating transaction", "from", transaction.FromAccountID, "to", transaction.ToAccountID, "amount", transaction.Amount)

	if transaction.Amount <= 0 {
		return nil, fmt.Errorf("%w: transaction amount must be positive", errs.ErrValidation)
	}

	if transaction.FromAccountID == transaction.ToAccountID {
		return nil, fmt.Errorf("%w: from and to accounts cannot be the same", errs.ErrValidation)
	}

	userID := ctxkeys.GetUserID(ctx)
	fromAccount, err := s.accountsClient.GetAccount(ctx, transaction.FromAccountID.String())
	if err != nil {
		return nil, fmt.Errorf("get from_account: %w", err)
	}
	s.logger.Infow("from_account fetched", "id", transaction.FromAccountID)

	if fromAccount == nil {
		return nil, fmt.Errorf("%w: from_account not found: %s", errs.ErrNotFound, transaction.FromAccountID)
	}

	// Validate ownership
	s.logger.Infow("checking ownership", "userID", userID, "account.UserID", fromAccount.UserID)
	if fromAccount.UserID != userID {
		return nil, fmt.Errorf("%w: account %s does not belong to user %s", errs.ErrUnauthorized, transaction.FromAccountID, userID)
	}

	s.logger.Infow("fetching to_account", "id", transaction.ToAccountID)
	toAccount, err := s.accountsClient.GetAccount(ctx, transaction.ToAccountID.String())
	if err != nil {
		return nil, fmt.Errorf("get to_account: %w", err)
	}
	s.logger.Infow("to_account fetched", "id", transaction.ToAccountID)

	if toAccount == nil {
		return nil, fmt.Errorf("%w: to_account not found: %s", errs.ErrNotFound, transaction.ToAccountID)
	}

	if fromAccount.Currency != toAccount.Currency {
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
