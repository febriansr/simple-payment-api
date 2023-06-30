package usecase

import (
	"github.com/febriansr/simple-payment-api/repository"
)

type LogoutUsecase interface {
	Logout(token string) error
}

type logoutUsecase struct {
	logoutRepository repository.LogoutRepository
}

func (l *logoutUsecase) Logout(token string) error {
	return l.logoutRepository.Logout(token)
}

func NewLogoutUsecase(logoutRepository repository.LogoutRepository) LogoutUsecase {
	return &logoutUsecase{
		logoutRepository: logoutRepository,
	}
}
