package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/httpclient/accounts"
	"go.uber.org/zap"
)

type Service interface {
	Hold(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error
	Deposit(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error
	Refund(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error
}

type implementation struct {
	accountClient *accounts.Client
	logger        *zap.SugaredLogger
}

func NewService(accountClient *accounts.Client, logger *zap.SugaredLogger) Service {
	return &implementation{
		accountClient: accountClient,
		logger:        logger.Named("account_service"),
	}
}

func (s *implementation) Hold(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error {
	if err := s.accountClient.Hold(ctx, accountID, transactionID, amount); err != nil {
		return fmt.Errorf("hold account: %w", err)
	}
	return nil
}

func (s *implementation) Deposit(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error {
	if err := s.accountClient.Deposit(ctx, accountID, transactionID, amount); err != nil {
		return fmt.Errorf("deposit account: %w", err)
	}
	return nil
}

func (s *implementation) Refund(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error {
	if err := s.accountClient.Refund(ctx, accountID, transactionID, amount); err != nil {
		return fmt.Errorf("refund account: %w", err)
	}
	return nil
}
