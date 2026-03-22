package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/account"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type GetAccountUsecase struct {
	usecase.BaseUsecase

	AccountService interface {
		GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
		GetByNumber(ctx context.Context, number string) (*models.Account, error)
		GetByNumberAndCurrency(ctx context.Context, number string, currency string) (*models.Account, error)
		GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error)
		GetByUserIDAndCurrency(ctx context.Context, userID string, currency string) (*models.Account, error)
		GetByTypeAndCurrency(ctx context.Context, accountType string, currency string) (*models.Account, error)
		GetAllByUserID(ctx context.Context, userID string) ([]models.Account, error)
	}
}

func NewGetAccountUsecase(base usecase.BaseUsecase, accountService account.Service) GetAccountUsecase {
	return GetAccountUsecase{
		BaseUsecase:    base,
		AccountService: accountService,
	}
}

func (uc *GetAccountUsecase) ExecuteByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	return uc.AccountService.GetByID(ctx, id)
}

func (uc *GetAccountUsecase) ExecuteByNumber(ctx context.Context, number string) (*models.Account, error) {
	return uc.AccountService.GetByNumber(ctx, number)
}

func (uc *GetAccountUsecase) ExecuteByNumberAndCurrency(ctx context.Context, number string, currency string) (*models.Account, error) {
	return uc.AccountService.GetByNumberAndCurrency(ctx, number, currency)
}

func (uc *GetAccountUsecase) ExecuteByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error) {
	return uc.AccountService.GetByUserID(ctx, userID)
}

func (uc *GetAccountUsecase) ExecuteByUserIDAndCurrency(ctx context.Context, userID string, currency string) (*models.Account, error) {
	return uc.AccountService.GetByUserIDAndCurrency(ctx, userID, currency)
}

func (uc *GetAccountUsecase) ExecuteByTypeAndCurrency(ctx context.Context, accountType string, currency string) (*models.Account, error) {
	return uc.AccountService.GetByTypeAndCurrency(ctx, accountType, currency)
}

func (uc *GetAccountUsecase) ExecuteAllByUserID(ctx context.Context, userID string) ([]models.Account, error) {
	return uc.AccountService.GetAllByUserID(ctx, userID)
}
