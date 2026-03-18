package users

import (
	"context"

	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type UpdateUserUsecase struct {
	usecase.BaseUsecase

	UsersService interface {
		Update(ctx context.Context, user models.User) (*models.User, error)
	}
}

func (uc *UpdateUserUsecase) Execute(ctx context.Context, user models.User) (*models.User, error) {
	uc.Logger.Infow("starting UpdateUserUsecase execution", "id", user.ID)

	txCtx, err := uc.Tx.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}

	defer uc.Tx.FinalizeTransaction(txCtx, &err)

	updatedUser, err := uc.UsersService.Update(txCtx, user)
	if err != nil {
		uc.Logger.Errorw("failed to update user", "id", user.ID, "error", err)
		return nil, err
	}
	uc.Logger.Infow("UpdateUserUsecase executed successfully", "id", updatedUser.ID)
	return updatedUser, nil
}
