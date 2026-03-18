package users

import (
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/users"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
)

type UsersDomain struct {
	CreateUsecase     CreateUsecase
	GetUserUsecase    GetUserUsecase
	UpdateUserUsecase UpdateUserUsecase
}

func NewDomain(baseusecase usecase.BaseUsecase, userService users.Service) UsersDomain {
	baseusecase.Logger = baseusecase.Logger.Named("users_domain")
	return UsersDomain{
		CreateUsecase: CreateUsecase{
			BaseUsecase:  baseusecase,
			UsersService: userService,
		},
		GetUserUsecase: GetUserUsecase{
			BaseUsecase:  baseusecase,
			UsersService: userService,
		},
		UpdateUserUsecase: UpdateUserUsecase{
			BaseUsecase:  baseusecase,
			UsersService: userService,
		},
	}
}
