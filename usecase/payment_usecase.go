package usecase

import (
	"github.com/febriansr/simple-payment-api/model/app_error"
	entity "github.com/febriansr/simple-payment-api/model/entity"
	"github.com/febriansr/simple-payment-api/repository"
)

type PaymentUsecase interface {
	PayTransaction(transaction entity.History) error
}

type paymentUsecase struct {
	paymentRepository repository.PaymentRepository
}

func (p *paymentUsecase) PayTransaction(transaction entity.History) error {
	if transaction.Amount < 0 {
		return app_error.InvalidError("invalid amount")
	}
	return p.paymentRepository.PayTransaction(transaction)
}

func NewPaymentUsecase(paymentRepository repository.PaymentRepository) PaymentUsecase {
	return &paymentUsecase{
		paymentRepository: paymentRepository,
	}
}
