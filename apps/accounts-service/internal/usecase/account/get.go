package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type GetAccountUsecase struct {
	usecase.BaseUsecase

	AccountService interface {
		GetByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
		GetByNumber(ctx context.Context, number string) (*models.Account, error)
		GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error)
	}
}

func (uc *GetAccountUsecase) ExecuteByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	return uc.AccountService.GetByID(ctx, id)
}

func (uc *GetAccountUsecase) ExecuteByNumber(ctx context.Context, number string) (*models.Account, error) {
	return uc.AccountService.GetByNumber(ctx, number)
}

func (uc *GetAccountUsecase) ExecuteByUserID(ctx context.Context, userID uuid.UUID) (*models.Account, error) {
	return uc.AccountService.GetByUserID(ctx, userID)
}
