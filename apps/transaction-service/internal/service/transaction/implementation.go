package transaction

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/ctxkeys"
	"github.com/w0ikid/yarmaq/pkg/errs"
	"github.com/w0ikid/yarmaq/pkg/exchange"
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
	exchange       exchange.Service
	logger         *zap.SugaredLogger
}

func NewService(repo TransactionRepo, accountsClient *accounts.Client, exchangeSvc exchange.Service, logger *zap.SugaredLogger) Service {
	return &implementation{
		repo:           repo,
		accountsClient: accountsClient,
		exchange:       exchangeSvc,
		logger:         logger.Named("transaction_service"),
	}
}

func (s *implementation) Create(ctx context.Context, transaction models.Transaction) (*models.Transaction, error) {
	s.logger.Infow("creating transaction", "to_account_number", transaction.ToAccountNumber, "amount", transaction.Amount, "currency", transaction.Currency)

	if transaction.Type == "" {
		transaction.Type = models.TransactionTypeTransfer
	}
	if !models.IsValidTransactionType(transaction.Type) {
		return nil, fmt.Errorf("%w: unsupported transaction type: %s", errs.ErrValidation, transaction.Type)
	}

	if transaction.Amount <= 0 {
		return nil, fmt.Errorf("%w: transaction amount must be positive", errs.ErrValidation)
	}

	userID := ctxkeys.GetUserID(ctx)
	if userID == "" {
		return nil, fmt.Errorf("%w: user is not authenticated", errs.ErrUnauthorized)
	}

	var fromAccount, toAccount *models.AccountResponse
	var err error

	switch transaction.Type {
	case models.TransactionTypeTransfer:
		if transaction.ToAccountNumber == "" {
			return nil, fmt.Errorf("%w: to_account_number is required", errs.ErrValidation)
		}

		fromAccount, err = s.accountsClient.GetAccountByUserIDAndCurrency(ctx, userID, transaction.Currency)
		if err != nil {
			return nil, fmt.Errorf("get from_account: %w", err)
		}
		if fromAccount == nil {
			return nil, fmt.Errorf("%w: from_account not found for user %s and currency %s", errs.ErrNotFound, userID, transaction.Currency)
		}

		toAccount, err = s.accountsClient.GetAccountByNumber(ctx, transaction.ToAccountNumber)
		if err != nil {
			return nil, fmt.Errorf("get to_account: %w", err)
		}
		if toAccount == nil {
			return nil, fmt.Errorf("%w: to_account not found for number %s", errs.ErrNotFound, transaction.ToAccountNumber)
		}
		if toAccount.UserID == userID {
			return nil, fmt.Errorf("%w: cannot transfer to your own account", errs.ErrValidation)
		}
	case models.TransactionTypeDeposit:
		fromAccount, err = s.accountsClient.GetSystemAccountByCurrency(ctx, transaction.Currency)
		if err != nil {
			return nil, fmt.Errorf("get system from_account: %w", err)
		}
		if fromAccount == nil {
			return nil, fmt.Errorf("%w: system account not found for currency %s", errs.ErrNotFound, transaction.Currency)
		}

		toAccount, err = s.accountsClient.GetAccountByUserIDAndCurrency(ctx, userID, transaction.Currency)
		if err != nil {
			return nil, fmt.Errorf("get user to_account: %w", err)
		}
		if toAccount == nil {
			return nil, fmt.Errorf("%w: user account not found for user %s and currency %s", errs.ErrNotFound, userID, transaction.Currency)
		}
	case models.TransactionTypeWithdrawal:
		fromAccount, err = s.accountsClient.GetAccountByUserIDAndCurrency(ctx, userID, transaction.Currency)
		if err != nil {
			return nil, fmt.Errorf("get user from_account: %w", err)
		}
		if fromAccount == nil {
			return nil, fmt.Errorf("%w: user account not found for user %s and currency %s", errs.ErrNotFound, userID, transaction.Currency)
		}

		toAccount, err = s.accountsClient.GetSystemAccountByCurrency(ctx, transaction.Currency)
		if err != nil {
			return nil, fmt.Errorf("get system to_account: %w", err)
		}
		if toAccount == nil {
			return nil, fmt.Errorf("%w: system account not found for currency %s", errs.ErrNotFound, transaction.Currency)
		}
	case models.TransactionTypeExchange:
		if transaction.TargetCurrency == nil {
			return nil, fmt.Errorf("%w: target_currency is required", errs.ErrValidation)
		}

		fromAccount, err = s.accountsClient.GetAccountByUserIDAndCurrency(ctx, userID, transaction.Currency)
		if err != nil {
			return nil, fmt.Errorf("get from_account: %w", err)
		}
		if fromAccount == nil {
			return nil, fmt.Errorf("%w: from_account not found for user %s and currency %s", errs.ErrNotFound, userID, transaction.Currency)
		}

		toAccount, err = s.accountsClient.GetAccountByUserIDAndCurrency(ctx, userID, *transaction.TargetCurrency)
		if err != nil {
			return nil, fmt.Errorf("get to_account: %w", err)
		}
		if toAccount == nil {
			return nil, fmt.Errorf("%w: to_account (currency: %s) not found for user %s", errs.ErrNotFound, *transaction.TargetCurrency, userID)
		}
	}

	transaction.FromAccountID = fromAccount.ID
	transaction.ToAccountID = toAccount.ID

	if transaction.FromAccountID == transaction.ToAccountID {
		return nil, fmt.Errorf("%w: from and to accounts cannot be the same", errs.ErrValidation)
	}

	// Calculate target amount and exchange rate
	targetAmount, rate, err := s.exchange.Convert(transaction.Amount, fromAccount.Currency, toAccount.Currency)
	if err != nil {
		return nil, fmt.Errorf("exchange conversion failed: %w", err)
	}

	if transaction.Amount > 0 && targetAmount == 0 {
		return nil, fmt.Errorf("%w: amount is too small for conversion from %s to %s", errs.ErrValidation, fromAccount.Currency, toAccount.Currency)
	}

	transaction.TargetAmount = &targetAmount
	transaction.TargetCurrency = &toAccount.Currency
	transaction.ExchangeRate = &rate
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
