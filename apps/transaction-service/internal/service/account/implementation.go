package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/pkg/httpclient/accounts"
	"go.uber.org/zap"
)

type Service interface {
	Debit(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error
	Credit(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error
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

func (s *implementation) Debit(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error {
	s.logger.Infow("debiting account", "accountID", accountID, "transactionID", transactionID, "amount", amount)
	err := s.accountClient.Debit(ctx, accountID, transactionID, amount)
	if err != nil {
		return fmt.Errorf("debit account: %w", err)
	}
	s.logger.Infow("account debited", "accountID", accountID, "transactionID", transactionID)
	return nil
}

func (s *implementation) Credit(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error {
	s.logger.Infow("crediting account", "accountID", accountID, "transactionID", transactionID, "amount", amount)
	err := s.accountClient.Credit(ctx, accountID, transactionID, amount)
	if err != nil {
		return fmt.Errorf("credit account: %w", err)
	}
	s.logger.Infow("account credited", "accountID", accountID, "transactionID", transactionID)
	return nil
}
