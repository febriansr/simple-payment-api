package usecase

import (
	"errors"
	"testing"

	entity "github.com/febriansr/simple-payment-api/model/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyTransaction = []entity.History{
	{
		MerchantCode: "Dummy Merchant Code",
		Amount:       20000.00,
	},
	{
		MerchantCode: "Dummy Merchant Code",
		Amount:       -30000.00,
	},
}

type paymentRepoMock struct {
	mock.Mock
}

func (p *paymentRepoMock) PayTransaction(transaction entity.History) error {
	args := p.Called(transaction)
	if args[0] != nil {
		return errors.New("Failed")
	}
	return nil
}

type PaymentUsecaseTestSuite struct {
	paymentRepoMock *paymentRepoMock
	suite.Suite
}

func (suite *PaymentUsecaseTestSuite) TestPayTransaction_Success() {
	paymentUsecase := NewPaymentUsecase(suite.paymentRepoMock)
	suite.paymentRepoMock.On("PayTransaction", dummyTransaction[0]).Return(nil)
	err := paymentUsecase.PayTransaction(dummyTransaction[0])
	assert.Nil(suite.T(), err)
}

func (suite *PaymentUsecaseTestSuite) TestPayTransaction_FailedRepo() {
	paymentUsecase := NewPaymentUsecase(suite.paymentRepoMock)
	suite.paymentRepoMock.On("PayTransaction", dummyTransaction[0]).Return(errors.New("failed"))
	err := paymentUsecase.PayTransaction(dummyTransaction[0])
	assert.NotNil(suite.T(), err)
}

func (suite *PaymentUsecaseTestSuite) TestPayTransaction_FailedInvalidAmount() {
	paymentUsecase := NewPaymentUsecase(suite.paymentRepoMock)
	suite.paymentRepoMock.On("PayTransaction", dummyTransaction[1]).Return(errors.New("failed"))
	err := paymentUsecase.PayTransaction(dummyTransaction[1])
	assert.NotNil(suite.T(), err)
}

func (suite *PaymentUsecaseTestSuite) SetupTest() {
	suite.paymentRepoMock = new(paymentRepoMock)
}

func TestPaymentUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentUsecaseTestSuite))
}
