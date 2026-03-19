package account

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/models"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, account models.Account) (*models.Account, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	GetByNumber(ctx context.Context, number string) (*models.Account, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error)
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
	s.logger.Infow("creating account", "user_id", account.UserID, "number", account.Number)
	
	seq, err := s.repo.NextSeq(ctx)
	if err != nil {
		s.logger.Errorw("failed to get next account number sequence", "error", err)
		return nil, err
	}

	// Number generation logic
	account.Number = generateAccountNumber(account.Currency, seq)
	
	return s.repo.Create(ctx, account)
}

func (s *implementation) GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *implementation) GetByNumber(ctx context.Context, number string) (*models.Account, error) {
	return s.repo.GetByNumber(ctx, number)
}

func (s *implementation) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *implementation) UpdateBalance(ctx context.Context, accountID uuid.UUID, amount int64, operationType string, referenceID *uuid.UUID) error {
	s.logger.Infow("updating balance", "account_id", accountID, "amount", amount, "op", operationType)

	// 1. Get account
	acc, err := s.repo.GetByID(ctx, accountID)
	if err != nil {
		return err
	}
	if acc == nil {
		return fmt.Errorf("account not found: %s", accountID)
	}

	// 2. Update balance
	acc.Balance += amount
	if acc.Balance < 0 {
		return fmt.Errorf("insufficient funds in account: %s", accountID)
	}

	_, err = s.repo.Update(ctx, *acc)
	if err != nil {
		return err
	}

	// 3. Create Ledger entry
	_, err = s.ledgerRepo.Create(ctx, models.Ledger{
		ID:            uuid.New(),
		AccountID:     accountID,
		Amount:        amount,
		OperationType: operationType,
		ReferenceID:   referenceID,
		CreatedAt:     time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
