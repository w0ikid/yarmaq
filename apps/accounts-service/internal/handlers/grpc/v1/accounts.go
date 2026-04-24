package v1

import (
	"context"

	"github.com/google/uuid"
	accountsv1 "github.com/w0ikid/yarmaq/pkg/gen/accounts/v1"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase/account"
	"github.com/w0ikid/yarmaq/pkg/models"
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

func (h *AccountsHandler) GetAccountByID(ctx context.Context, req *accountsv1.GetAccountByIDRequest) (*accountsv1.GetAccountResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	acc, err := h.domain.GetAccountUsecase.ExecuteByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get account: %v", err)
	}
	if acc == nil {
		return nil, status.Errorf(codes.NotFound, "account not found")
	}

	return h.mapToProto(acc), nil
}

func (h *AccountsHandler) GetAccountByNumber(ctx context.Context, req *accountsv1.GetAccountByNumberRequest) (*accountsv1.GetAccountResponse, error) {
	acc, err := h.domain.GetAccountUsecase.ExecuteByNumber(ctx, req.Number)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get account: %v", err)
	}
	if acc == nil {
		return nil, status.Errorf(codes.NotFound, "account not found")
	}

	return h.mapToProto(acc), nil
}

func (h *AccountsHandler) GetAccountByUserIDAndCurrency(ctx context.Context, req *accountsv1.GetAccountByUserIDAndCurrencyRequest) (*accountsv1.GetAccountResponse, error) {
	acc, err := h.domain.GetAccountUsecase.ExecuteByUserIDAndCurrency(ctx, req.UserId, req.Currency)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get account: %v", err)
	}
	if acc == nil {
		return nil, status.Errorf(codes.NotFound, "account not found")
	}

	return h.mapToProto(acc), nil
}

func (h *AccountsHandler) GetSystemAccountByCurrency(ctx context.Context, req *accountsv1.GetSystemAccountByCurrencyRequest) (*accountsv1.GetAccountResponse, error) {
	// AccountTypeSystem is a known value
	acc, err := h.domain.GetAccountUsecase.ExecuteByTypeAndCurrency(ctx, "SYSTEM", req.Currency)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get system account: %v", err)
	}
	if acc == nil {
		return nil, status.Errorf(codes.NotFound, "system account not found")
	}

	return h.mapToProto(acc), nil
}

func (h *AccountsHandler) mapToProto(acc *models.Account) *accountsv1.GetAccountResponse {
	userID := ""
	if acc.UserID != nil {
		userID = *acc.UserID
	}
	return &accountsv1.GetAccountResponse{
		Id:       acc.ID.String(),
		UserId:   userID,
		Type:     acc.Type,
		Number:   acc.Number,
		Balance:  acc.Balance,
		Currency: acc.Currency,
		Status:   acc.Status,
	}
}
