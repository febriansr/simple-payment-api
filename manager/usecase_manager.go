package manager

import (
	"github.com/febriansr/simple-payment-api/usecase"
	"github.com/febriansr/simple-payment-api/utils/authenticator"
)

type UsecaseManager interface {
	LoginUsecase() usecase.LoginUsecase
	LogoutUsecase() usecase.LogoutUsecase
	PaymentUsecase() usecase.PaymentUsecase
}

type usecaseManager struct {
	repositoryManager RepositoryManager
	authenticator     authenticator.AccessToken
}

func (u *usecaseManager) LoginUsecase() usecase.LoginUsecase {
	return usecase.NewLoginUsecase(u.repositoryManager.LoginRepository(), u.authenticator)
}

func (u *usecaseManager) LogoutUsecase() usecase.LogoutUsecase {
	return usecase.NewLogoutUsecase(u.repositoryManager.LogoutRepository())
}

func (u *usecaseManager) PaymentUsecase() usecase.PaymentUsecase {
	return usecase.NewPaymentUsecase(u.repositoryManager.PaymentRepository())
}

func NewUsecaseManager(r RepositoryManager, a authenticator.AccessToken) UsecaseManager {
	return &usecaseManager{
		repositoryManager: r,
		authenticator:     a,
	}
}
