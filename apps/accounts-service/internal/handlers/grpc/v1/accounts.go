package v1

import (
	"context"

	"github.com/google/uuid"
	accountsv1 "github.com/w0ikid/yarmaq/pkg/gen/accounts/v1"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/account"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountsHandler struct {
	accountsv1.UnimplementedAccountsServiceServer
	domain account.AccountDomain
	logger *zap.SugaredLogger
}

func NewAccountsHandler(domain account.AccountDomain, logger *zap.SugaredLogger) *AccountsHandler {
	return &AccountsHandler{
		domain: domain,
		logger: logger.Named("grpc_v1_accounts"),
	}
}

func (h *AccountsHandler) UpdateBalance(ctx context.Context, req *accountsv1.UpdateBalanceRequest) (*accountsv1.UpdateBalanceResponse, error) {
	accountID, err := uuid.Parse(req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account_id: %v", err)
	}

	var refID *uuid.UUID
	if req.ReferenceId != "" {
		id, err := uuid.Parse(req.ReferenceId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid reference_id: %v", err)
		}
		refID = &id
	}

	err = h.domain.UpdateBalanceUsecase.Execute(ctx, accountID, req.Amount, req.OperationType, refID)
	if err != nil {
		h.logger.Errorw("failed to update balance", "account_id", accountID, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update balance: %v", err)
	}

	return &accountsv1.UpdateBalanceResponse{}, nil
}
