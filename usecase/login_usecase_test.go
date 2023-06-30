package usecase

import (
	"errors"
	"testing"

	entity "github.com/febriansr/simple-payment-api/model/entity"
	"github.com/febriansr/simple-payment-api/utils/authenticator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyCustomer = &entity.Customer{
	Username: "dummyUsername",
	Password: "dummyPassword",
}

var dummyTokenDetails = []authenticator.TokenDetails{
	{
		AccessToken: "Dummy Access Token",
		AccessUuid:  "Dummy Access Uuid",
		AtExpires:   5,
	},
}

type loginRepoMock struct {
	mock.Mock
}

func (l *loginRepoMock) FindCustomer(customer entity.Customer) error {
	args := l.Called(customer)
	if args[0] != nil {
		return errors.New("Failed")
	}
	return nil
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
	return args.Get(0).(authenticator.AccessDetails), nil
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

type LoginUsecaseTestSuite struct {
	loginRepoMock *loginRepoMock
	authMock      *authMock
	suite.Suite
}

func (suite *LoginUsecaseTestSuite) TestLogin_Success() {
	loginUsecase := NewLoginUsecase(suite.loginRepoMock, suite.authMock)
	suite.loginRepoMock.On("FindCustomer", *dummyCustomer).Return(nil)
	suite.authMock.On("CreateAccessToken", dummyCustomer).Return(dummyTokenDetails[0], nil)
	suite.authMock.On("StoreAccessToken", dummyCustomer.Username, dummyTokenDetails[0]).Return(nil)
	tokenDetails, err := loginUsecase.Login(*dummyCustomer)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyTokenDetails[0].AccessToken, tokenDetails)
}

func (suite *LoginUsecaseTestSuite) TestLogin_FailedFindCustomer() {
	loginUsecase := NewLoginUsecase(suite.loginRepoMock, suite.authMock)
	suite.loginRepoMock.On("FindCustomer", *dummyCustomer).Return(errors.New("Failed"))
	tokenDetails, err := loginUsecase.Login(*dummyCustomer)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", tokenDetails)
}

func (suite *LoginUsecaseTestSuite) TestLogin_FailedCreateAccessToken() {
	loginUsecase := NewLoginUsecase(suite.loginRepoMock, suite.authMock)
	suite.loginRepoMock.On("FindCustomer", *dummyCustomer).Return(nil)
	suite.authMock.On("CreateAccessToken", dummyCustomer).Return(authenticator.TokenDetails{}, errors.New("Failed"))
	tokenDetails, err := loginUsecase.Login(*dummyCustomer)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", tokenDetails)
}

func (suite *LoginUsecaseTestSuite) TestLogin_FailedStoreAccessToken() {
	loginUsecase := NewLoginUsecase(suite.loginRepoMock, suite.authMock)
	suite.loginRepoMock.On("FindCustomer", *dummyCustomer).Return(nil)
	suite.authMock.On("CreateAccessToken", dummyCustomer).Return(dummyTokenDetails[0], nil)
	suite.authMock.On("StoreAccessToken", dummyCustomer.Username, dummyTokenDetails[0]).Return(errors.New("Failed"))
	tokenDetails, err := loginUsecase.Login(*dummyCustomer)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "", tokenDetails)
}

func (suite *LoginUsecaseTestSuite) SetupTest() {
	suite.loginRepoMock = new(loginRepoMock)
	suite.authMock = new(authMock)
}

func TestLoginUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(LoginUsecaseTestSuite))
}
