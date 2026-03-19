package account

import (
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/service/account"
	"github.com/w0ikid/yarmaq/apps/accounts-service/internal/usecase"
)

type AccountDomain struct {
	CreateUsecase        CreateAccountUsecase
	GetAccountUsecase    GetAccountUsecase
	UpdateBalanceUsecase UpdateBalanceUsecase
}

func NewDomain(baseusecase usecase.BaseUsecase, accountService account.Service) AccountDomain {
	baseusecase.Logger = baseusecase.Logger.Named("account_domain")
	return AccountDomain{
		CreateUsecase: CreateAccountUsecase{
			BaseUsecase:    baseusecase,
			AccountService: accountService,
		},
		GetAccountUsecase: GetAccountUsecase{
			BaseUsecase:    baseusecase,
			AccountService: accountService,
		},
		UpdateBalanceUsecase: UpdateBalanceUsecase{
			BaseUsecase:    baseusecase,
			AccountService: accountService,
		},
	}
}
