package manager

import (
	"github.com/febriansr/simple-payment-api/config"
	"github.com/febriansr/simple-payment-api/repository"
	"github.com/febriansr/simple-payment-api/utils/authenticator"
)

type RepositoryManager interface {
	LoginRepository() repository.LoginRepository
	LogoutRepository() repository.LogoutRepository
	PaymentRepository() repository.PaymentRepository
}

type repositoryManager struct {
	config        config.JsonFileConfig
	authenticator authenticator.AccessToken
}

func (r *repositoryManager) LoginRepository() repository.LoginRepository {
	return repository.NewLoginRepository(r.config)
}

func (r *repositoryManager) LogoutRepository() repository.LogoutRepository {
	return repository.NewLogoutRepository(r.authenticator)
}

func (r *repositoryManager) PaymentRepository() repository.PaymentRepository {
	return repository.NewPaymentRepository(r.config)
}

func NewRepositoryManager(config config.JsonFileConfig, authenticator authenticator.AccessToken) RepositoryManager {
	return &repositoryManager{
		config:        config,
		authenticator: authenticator,
	}
}
