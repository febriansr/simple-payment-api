package repository

import (
	"github.com/febriansr/simple-payment-api/config"
	"github.com/febriansr/simple-payment-api/model/app_error"
	entity "github.com/febriansr/simple-payment-api/model/entity"
	"github.com/febriansr/simple-payment-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type LoginRepository interface {
	FindCustomer(iCustomer entity.Customer) error
}

type loginRepository struct {
	config config.JsonFileConfig
}

func (l *loginRepository) FindCustomer(iCustomer entity.Customer) error {
	var customers []entity.Customer
	err := utils.ReadParseJSON(l.config.Customer, &customers)
	if err != nil {
		return app_error.InternalServerError("Failed to read and parse customer data: " + err.Error())
	}

	for _, customer := range customers {
		if customer.Username == iCustomer.Username {
			err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(iCustomer.Password))
			if err != nil {
				return app_error.Unauthorized("Invalid credentials: " + err.Error())
			}
			return nil
		}
	}

	return app_error.DataNotFound("user not found")
}

func NewLoginRepository(config config.JsonFileConfig) LoginRepository {
	return &loginRepository{
		config: config,
	}
}
