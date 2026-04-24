package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	accountsv1 "github.com/w0ikid/yarmaq/pkg/gen/accounts/v1"
	"go.uber.org/zap"
)

type Service interface {
	Hold(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error
	Deposit(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error
	Refund(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error
}

type implementation struct {
	accountClient accountsv1.AccountsServiceClient
	logger        *zap.SugaredLogger
}

func NewService(accountClient accountsv1.AccountsServiceClient, logger *zap.SugaredLogger) Service {
	return &implementation{
		accountClient: accountClient,
		logger:        logger.Named("account_service"),
	}
}

func (s *implementation) Hold(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error {
	_, err := s.accountClient.UpdateBalance(ctx, &accountsv1.UpdateBalanceRequest{
		AccountId:     accountID,
		Amount:        -amount, // Decrease balance
		OperationType: "HOLD",
		ReferenceId:   transactionID.String(),
	})
	if err != nil {
		return fmt.Errorf("hold account (gRPC): %w", err)
	}
	return nil
}

func (s *implementation) Deposit(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error {
	_, err := s.accountClient.UpdateBalance(ctx, &accountsv1.UpdateBalanceRequest{
		AccountId:     accountID,
		Amount:        amount, // Increase balance
		OperationType: "DEPOSIT",
		ReferenceId:   transactionID.String(),
	})
	if err != nil {
		return fmt.Errorf("deposit account (gRPC): %w", err)
	}
	return nil
}

func (s *implementation) Refund(ctx context.Context, accountID string, transactionID uuid.UUID, amount int64) error {
	_, err := s.accountClient.UpdateBalance(ctx, &accountsv1.UpdateBalanceRequest{
		AccountId:     accountID,
		Amount:        amount, // Increase balance (refund)
		OperationType: "REFUND",
		ReferenceId:   transactionID.String(),
	})
	if err != nil {
		return fmt.Errorf("refund account (gRPC): %w", err)
	}
	return nil
}
