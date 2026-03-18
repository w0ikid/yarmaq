package users

import (
	"context"

	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
	"github.com/w0ikid/yarmaq/pkg/models"
)

type GetUserUsecase struct {
	usecase.BaseUsecase

	UsersService interface {
		GetByZitadelID(ctx context.Context, zitadelID string) (*models.User, error)
	}

}

func (uc *GetUserUsecase) Execute(ctx context.Context, zitadelID string) (*models.User, error) {
    uc.Logger.Infow("starting GetUserUsecase execution", "zitadelID", zitadelID)
    
	user, err := uc.UsersService.GetByZitadelID(ctx, zitadelID)
	if err != nil {
		uc.Logger.Errorw("error getting user by Zitadel ID", "zitadelID", zitadelID, "error", err)
		return nil, err
	}
	if user == nil {
		uc.Logger.Infow("user not found with Zitadel ID", "zitadelID", zitadelID)
		return nil, nil
	}
	
	uc.Logger.Infow("successfully got user by Zitadel ID", "zitadelID", zitadelID, "userID", user.ID)
	return user, nil
}
