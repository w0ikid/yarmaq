package users

import (
	"context"

	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type CreateUsecase struct {
	usecase.BaseUsecase

	UsersService interface {
		Create(ctx context.Context, user models.User) (*models.User, error)
		AssignRole(ctx context.Context, zitadelID string, roles []string) error
	}
}

func (uc *CreateUsecase) Execute(ctx context.Context, user models.User) (*models.User, error) {
	uc.Logger.Infow("starting CreateUsecase execution", "email", user.Email)

	txCtx, err := uc.Tx.StartTransaction(ctx)
	if err != nil {
		return nil, err
	}

	defer uc.Tx.FinalizeTransaction(txCtx, &err)

	createdUser, err := uc.UsersService.Create(txCtx, user)
	if err != nil {
		uc.Logger.Errorw("failed to create user", "email", user.Email, "error", err)
		return nil, err
	}

	assignErr := uc.UsersService.AssignRole(txCtx, createdUser.ZitadelUserID, user.Roles)
	if assignErr != nil {
		uc.Logger.Errorw("failed to assign role to user", "zitadelID", createdUser.ZitadelUserID, "role", user.Roles, "error", assignErr)
		return nil, assignErr
	}

	uc.Logger.Infow("CreateUsecase executed successfully", "id", createdUser.ID)
	return createdUser, nil
}
