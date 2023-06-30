package repository

import (
	"time"

	"github.com/febriansr/simple-payment-api/config"
	"github.com/febriansr/simple-payment-api/model/app_error"
	entity "github.com/febriansr/simple-payment-api/model/entity"
	"github.com/febriansr/simple-payment-api/utils"
	"github.com/google/uuid"
)

type PaymentRepository interface {
	PayTransaction(transaction entity.History) error
}

type paymentRepository struct {
	config config.JsonFileConfig
}

func (p *paymentRepository) PayTransaction(transaction entity.History) error {
	var customers []entity.Customer
	var merchants []entity.Merchant
	var histories []entity.History
	err := utils.ReadParseJSON(p.config.Customer, &customers)
	if err != nil {
		return app_error.InternalServerError("Failed to read and parse customer data: " + err.Error())
	}

	err = utils.ReadParseJSON(p.config.Merchant, &merchants)
	if err != nil {
		return app_error.InternalServerError("Failed to read and parse merchant data: " + err.Error())
	}

	err = utils.ReadParseJSON(p.config.History, &histories)
	if err != nil {
		return app_error.InternalServerError("Failed to read and parse history data: " + err.Error())
	}

	isCustomer := false
	isMerchant := false

	for i, customer := range customers {
		if customer.Username == transaction.CustomerUsername {
			isCustomer = true
			for _, merchant := range merchants {
				if merchant.MerchantCode == transaction.MerchantCode {
					isMerchant = true
					if customer.Balance >= transaction.Amount {
						transaction.Date = time.Now()
						transaction.TransactionId = uuid.New().String()
						customers[i].Balance -= transaction.Amount
						break
					} else {
						return app_error.InvalidError("Balance insufficient")
					}
				}
			}
		}
	}

	if !isCustomer {
		return app_error.InvalidError("Invalid username")
	}

	if !isMerchant {
		return app_error.InvalidError("Invalid merchant code")
	}

	err = utils.WriteJSON(p.config.Customer, customers)
	if err != nil {
		return app_error.InternalServerError("Failed to write updated customer data to file: " + err.Error())
	}

	err = utils.WriteJSON(p.config.History, histories)
	if err != nil {
		return app_error.InternalServerError("Failed to write updated history data to file: " + err.Error())
	}

	return nil
}

func NewPaymentRepository(config config.JsonFileConfig) PaymentRepository {
	return &paymentRepository{
		config: config,
	}
}
