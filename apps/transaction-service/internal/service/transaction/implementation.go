package transaction

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/ctxkeys"
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

	userID := ctxkeys.GetUserID(ctx)

	s.logger.Infow("fetching from_account", "id", transaction.FromAccountID)
	fromAccount, err := s.accountsClient.GetAccount(ctx, transaction.FromAccountID.String())
	if err != nil {
		return nil, fmt.Errorf("get from_account: %w", err)
	}
	s.logger.Infow("from_account fetched", "id", transaction.FromAccountID)

	if fromAccount == nil {
		return nil, fmt.Errorf("from_account not found: %s", transaction.FromAccountID)
	}

	// Validate ownership
	if fromAccount.UserID != userID {
		return nil, fmt.Errorf("unauthorized: account %s does not belong to user %s", transaction.FromAccountID, userID)
	}

	s.logger.Infow("fetching to_account", "id", transaction.ToAccountID)
	toAccount, err := s.accountsClient.GetAccount(ctx, transaction.ToAccountID.String())
	if err != nil {
		return nil, fmt.Errorf("get to_account: %w", err)
	}
	s.logger.Infow("to_account fetched", "id", transaction.ToAccountID)

	if toAccount == nil {
		return nil, fmt.Errorf("to_account not found: %s", transaction.ToAccountID)
	}

	if fromAccount.Currency != toAccount.Currency {
		return nil, fmt.Errorf("accounts have different currencies: %s vs %s", fromAccount.Currency, toAccount.Currency)
	}

	if transaction.Amount <= 0 {
		return nil, fmt.Errorf("transaction amount must be positive")
	}

	if transaction.FromAccountID == transaction.ToAccountID {
		return nil, fmt.Errorf("from and to accounts cannot be the same")
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
		return fmt.Errorf("transaction not found: %s", id)
	}

	tx.Status = status
	_, err = s.repo.Update(ctx, *tx)
	return err
}
