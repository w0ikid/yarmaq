package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/errs"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, account models.Account) (*models.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	GetByNumber(ctx context.Context, number string) (*models.Account, error)
	GetByNumberAndCurrency(ctx context.Context, number string, currency string) (*models.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error)
	GetByUserIDAndCurrency(ctx context.Context, userID string, currency string) (*models.Account, error)
	GetByTypeAndCurrency(ctx context.Context, accountType string, currency string) (*models.Account, error)
	GetAllByUserID(ctx context.Context, userID string) ([]models.Account, error)
	UpdateBalance(ctx context.Context, accountID uuid.UUID, amount int64, operationType string, referenceID *uuid.UUID) error
}

type implementation struct {
	repo       AccountRepo
	ledgerRepo LedgerRepo
	outboxRepo OutboxRepo
	logger     *zap.SugaredLogger
}

func NewService(repo AccountRepo, ledgerRepo LedgerRepo, outboxRepo OutboxRepo, logger *zap.SugaredLogger) Service {
	return &implementation{
		repo:       repo,
		ledgerRepo: ledgerRepo,
		outboxRepo: outboxRepo,
		logger:     logger.Named("account_service"),
	}
}

func (s *implementation) Create(ctx context.Context, account models.Account) (*models.Account, error) {
	if !models.IsValidCurrency(account.Currency) {
		return nil, fmt.Errorf("%w: unsupported currency: %s", errs.ErrValidation, account.Currency)
	}

	seq, err := s.repo.NextSeq(ctx)
	if err != nil {
		s.logger.Errorw("failed to get next seq", "error", err)
		return nil, err
	}

	account.Number, err = generateAccountNumber(account.Currency, seq)
	if err != nil {
		return nil, err
	}

	s.logger.Infow("creating account", "user_id", account.UserID, "number", account.Number)
	return s.repo.Create(ctx, account)
}

func (s *implementation) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *implementation) GetByNumber(ctx context.Context, number string) (*models.Account, error) {
	return s.repo.GetByNumber(ctx, number)
}

func (s *implementation) GetByNumberAndCurrency(ctx context.Context, number string, currency string) (*models.Account, error) {
	return s.repo.GetByNumberAndCurrency(ctx, number, currency)
}

func (s *implementation) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *implementation) GetByUserIDAndCurrency(ctx context.Context, userID string, currency string) (*models.Account, error) {
	return s.repo.GetByUserIDAndCurrency(ctx, userID, currency)
}

func (s *implementation) GetByTypeAndCurrency(ctx context.Context, accountType string, currency string) (*models.Account, error) {
	return s.repo.GetByTypeAndCurrency(ctx, accountType, currency)
}

func (s *implementation) GetAllByUserID(ctx context.Context, userID string) ([]models.Account, error) {
	return s.repo.GetAllByUserID(ctx, userID)
}

func (s *implementation) UpdateBalance(ctx context.Context, accountID uuid.UUID, amount int64, operationType string, referenceID *uuid.UUID) error {
	s.logger.Infow("updating balance", "account_id", accountID, "amount", amount, "op", operationType)

	acc, err := s.repo.GetByID(ctx, accountID)
	if err != nil {
		s.logger.Errorw("failed to get account", "error", err)
		return err
	}
	if acc == nil {
		return fmt.Errorf("%w: account not found: %s", errs.ErrNotFound, accountID)
	}

	acc.Balance += amount
	if acc.Balance < 0 {
		return fmt.Errorf("%w: insufficient funds: %s", errs.ErrValidation, accountID)
	}

	if _, err = s.repo.Update(ctx, *acc); err != nil {
		s.logger.Errorw("failed to update account", "error", err)
		return err
	}

	_, err = s.ledgerRepo.Create(ctx, models.Ledger{
		AccountID:     accountID,
		Amount:        amount,
		OperationType: operationType,
		ReferenceID:   referenceID,
	})
	if err != nil {
		s.logger.Errorw("failed to create ledger", "error", err)
	}
	return err
}
