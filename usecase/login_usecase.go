package usecase

import (
	entity "github.com/febriansr/simple-payment-api/model/entity"
	"github.com/febriansr/simple-payment-api/repository"
	"github.com/febriansr/simple-payment-api/utils/authenticator"
)

type LoginUsecase interface {
	Login(customer entity.Customer) (token string, err error)
}

type loginUsecase struct {
	loginRepository repository.LoginRepository
	authenticator   authenticator.AccessToken
}

func (l *loginUsecase) Login(customer entity.Customer) (token string, err error) {
	res := l.loginRepository.FindCustomer(customer)

	if res == nil {
		tokenDetails, err := l.authenticator.CreateAccessToken(&customer)
		if err != nil {
			return "", err
		}
		err = l.authenticator.StoreAccessToken(customer.Username, tokenDetails)
		if err != nil {
			return "", err
		}
		return tokenDetails.AccessToken, nil
	} else {
		return "", res
	}
}

func NewLoginUsecase(loginRepository repository.LoginRepository, authenticator authenticator.AccessToken) LoginUsecase {
	return &loginUsecase{
		loginRepository: loginRepository,
		authenticator:   authenticator,
	}
}
