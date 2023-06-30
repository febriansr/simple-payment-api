package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/febriansr/simple-payment-api/model/dto/res"
	entity "github.com/febriansr/simple-payment-api/model/entity"
	"github.com/febriansr/simple-payment-api/utils/authenticator"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyTransaction = []entity.History{
	{
		MerchantCode: "Dummy Merchant Code",
		Amount:       20000.00,
	},
}

var dummyAccessDetails = []authenticator.AccessDetails{
	{
		AccessUuid: "Dummy Access Uuid",
		Username:   "Dummy Username",
	},
}

type paymentUsecaseMock struct {
	mock.Mock
}

func (p *paymentUsecaseMock) PayTransaction(transaction entity.History) error {
	args := p.Called(transaction)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

type bindAuthHeaderMock struct {
	mock.Mock
}

func (b *bindAuthHeaderMock) BindAuthHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Empty token")
	}
	return dummyTokenDetails[0].AccessToken, nil
}

type authMock struct {
	mock.Mock
}

func (a *authMock) CreateAccessToken(customer *entity.Customer) (authenticator.TokenDetails, error) {
	args := a.Called(customer)
	if args == nil {
		return authenticator.TokenDetails{}, errors.New("Failed")
	}
	return args.Get(0).(authenticator.TokenDetails), args.Error(1)
}

func (a *authMock) VerifyAccessToken(tokenString string) (authenticator.AccessDetails, error) {
	args := a.Called(tokenString)
	if args == nil {
		return authenticator.AccessDetails{}, errors.New("Failed")
	}
	return args.Get(0).(authenticator.AccessDetails), args.Error(1)
}

func (a *authMock) StoreAccessToken(username string, tokenDetails authenticator.TokenDetails) error {
	args := a.Called(username, tokenDetails)
	if args[0] != nil {
		return errors.New("Failed")
	}
	return nil
}

func (a *authMock) FetchAccessToken(accessDetails authenticator.AccessDetails) error {
	args := a.Called(accessDetails)
	if args[0] != nil {
		return errors.New("Failed")
	}
	return nil
}

func (a *authMock) DeleteAccessToken(accessUuid string) error {
	args := a.Called(accessUuid)
	if args[0] != nil {
		return errors.New("Failed")
	}
	return nil
}

type middlewareMock struct {
	mock.Mock
}

func (m *middlewareMock) RequireToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

type PaymentControllerTestSuite struct {
	suite.Suite
	routerMock         *gin.Engine
	routerGroupMock    *gin.RouterGroup
	usecaseMock        *paymentUsecaseMock
	bindAuthHeaderMock *bindAuthHeaderMock
	authMock           *authMock
	middlewareMock     *middlewareMock
}

func (suite *PaymentControllerTestSuite) TestPayTransaction_Success() {
	transaction := dummyTransaction[0]
	paymentController := NewPaymentController(suite.routerGroupMock, suite.usecaseMock, suite.authMock, suite.middlewareMock)
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(transaction)
	request, _ := http.NewRequest(http.MethodPost, "/v1/payment", bytes.NewBuffer(reqBody))
	request.Header.Set("Authorization", dummyTokenDetails[0].AccessToken)
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = request
	suite.bindAuthHeaderMock.On("BindAuthHeader", ctx).Return(dummyTokenDetails[0].AccessToken, nil)
	suite.authMock.On("VerifyAccessToken", dummyTokenDetails[0].AccessToken).Return(dummyAccessDetails[0], nil)
	transaction.CustomerUsername = dummyAccessDetails[0].Username
	suite.usecaseMock.On("PayTransaction", transaction).Return(nil)

	paymentController.PaymentHandler(ctx)

	var response res.ApiResponse
	json.Unmarshal(r.Body.Bytes(), &response)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
}

func (suite *PaymentControllerTestSuite) TestPayTransaction_FailedBindJSON() {
	paymentController := NewPaymentController(suite.routerGroupMock, suite.usecaseMock, suite.authMock, suite.middlewareMock)
	r := httptest.NewRecorder()
	reqBody := []byte(`{1}`)
	request, _ := http.NewRequest(http.MethodPost, "/v1/payment", bytes.NewBuffer(reqBody))
	request.Header.Set("Authorization", dummyTokenDetails[0].AccessToken)
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = request

	paymentController.PaymentHandler(ctx)

	var response res.ApiResponse
	json.Unmarshal(r.Body.Bytes(), &response)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *PaymentControllerTestSuite) TestPayTransaction_FailedBindAuthHeader() {
	transaction := dummyTransaction[0]
	paymentController := NewPaymentController(suite.routerGroupMock, suite.usecaseMock, suite.authMock, suite.middlewareMock)
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(transaction)
	request, _ := http.NewRequest(http.MethodPost, "/v1/payment", bytes.NewBuffer(reqBody))
	request.Header.Set("Authorization", "Bearer")
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = request
	suite.bindAuthHeaderMock.On("BindAuthHeader", ctx).Return("", errors.New("Failed"))

	paymentController.PaymentHandler(ctx)

	var response res.ApiResponse
	json.Unmarshal(r.Body.Bytes(), &response)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *PaymentControllerTestSuite) TestPayTransaction_FailedVerifyAccessToken() {
	transaction := dummyTransaction[0]
	paymentController := NewPaymentController(suite.routerGroupMock, suite.usecaseMock, suite.authMock, suite.middlewareMock)
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(transaction)
	request, _ := http.NewRequest(http.MethodPost, "/v1/payment", bytes.NewBuffer(reqBody))
	request.Header.Set("Authorization", dummyTokenDetails[0].AccessToken)
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = request
	suite.bindAuthHeaderMock.On("BindAuthHeader", ctx).Return(dummyTokenDetails[0].AccessToken, nil)
	suite.authMock.On("VerifyAccessToken", dummyTokenDetails[0].AccessToken).Return(authenticator.AccessDetails{}, errors.New("Failed"))

	paymentController.PaymentHandler(ctx)

	var response res.ApiResponse
	json.Unmarshal(r.Body.Bytes(), &response)

	assert.Equal(suite.T(), "XX", response.Status.Code)
}

func (suite *PaymentControllerTestSuite) TestPayTransaction_FailedUsecase() {
	transaction := dummyTransaction[0]
	paymentController := NewPaymentController(suite.routerGroupMock, suite.usecaseMock, suite.authMock, suite.middlewareMock)
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(transaction)
	request, _ := http.NewRequest(http.MethodPost, "/v1/payment", bytes.NewBuffer(reqBody))
	request.Header.Set("Authorization", dummyTokenDetails[0].AccessToken)
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = request
	suite.bindAuthHeaderMock.On("BindAuthHeader", ctx).Return(dummyTokenDetails[0].AccessToken, nil)
	suite.authMock.On("VerifyAccessToken", dummyTokenDetails[0].AccessToken).Return(dummyAccessDetails[0], nil)
	transaction.CustomerUsername = dummyAccessDetails[0].Username
	suite.usecaseMock.On("PayTransaction", transaction).Return(errors.New("Failed"))
	paymentController.PaymentHandler(ctx)

	var response res.ApiResponse
	json.Unmarshal(r.Body.Bytes(), &response)

	assert.Equal(suite.T(), "XX", response.Status.Code)
}

func (suite *PaymentControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.routerGroupMock = suite.routerMock.Group("/v1")
	suite.usecaseMock = new(paymentUsecaseMock)
	suite.bindAuthHeaderMock = new(bindAuthHeaderMock)
	suite.authMock = new(authMock)
	suite.middlewareMock = new(middlewareMock)
}

func TestPaymentControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentControllerTestSuite))
}
