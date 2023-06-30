package repository

import (
	"log"

	"github.com/febriansr/simple-payment-api/utils/authenticator"
)

type LogoutRepository interface {
	Logout(token string) error
}

type logoutRepository struct {
	authenticator authenticator.AccessToken
}

func (a *logoutRepository) Logout(token string) error {
	log.Print(token)
	accountDetails, err := a.authenticator.VerifyAccessToken(token)
	if err != nil {
		return err
	}
	err = a.authenticator.DeleteAccessToken(accountDetails.AccessUuid)
	if err != nil {
		return err
	}
	return nil
}

func NewLogoutRepository(authenticator authenticator.AccessToken) LogoutRepository {
	return &logoutRepository{
		authenticator: authenticator,
	}
}
