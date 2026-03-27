package account

import (
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/account"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/outbox"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/users"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
)

type AccountDomain struct {
	CreateUsecase        CreateAccountUsecase
	GetAccountUsecase    GetAccountUsecase
	UpdateBalanceUsecase UpdateBalanceUsecase
}

func NewDomain(baseusecase usecase.BaseUsecase, accountService account.Service, outbox outbox.Service, usersService users.Service) AccountDomain {
	baseusecase.Logger = baseusecase.Logger.Named("account_domain")
	return AccountDomain{
		CreateUsecase:        NewCreateAccountUsecase(baseusecase, accountService, outbox, usersService),
		GetAccountUsecase:    NewGetAccountUsecase(baseusecase, accountService),
		UpdateBalanceUsecase: NewUpdateBalanceUsecase(baseusecase, accountService, outbox),
	}
}
