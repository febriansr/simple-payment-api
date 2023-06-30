package repository

import (
	"errors"
	"testing"

	entity "github.com/febriansr/simple-payment-api/model/entity"
	"github.com/febriansr/simple-payment-api/utils/authenticator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var dummyAccessDetails = []authenticator.AccessDetails{
	{
		AccessUuid: "Dummy Access Uuid",
		Username:   "Dummy Username",
	},
}

var dummyTokenDetails = []authenticator.TokenDetails{
	{
		AccessToken: "Dummy Access Token",
		AccessUuid:  "Dummy Access Uuid",
		AtExpires:   5,
	},
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

type LogoutRepoTestSuite struct {
	authMock *authMock
	suite.Suite
}

func (suite *LogoutRepoTestSuite) TestLogout_Success() {
	logoutRepo := NewLogoutRepository(suite.authMock)
	suite.authMock.On("VerifyAccessToken", dummyTokenDetails[0].AccessToken).Return(dummyAccessDetails[0], nil)
	suite.authMock.On("DeleteAccessToken", dummyAccessDetails[0].AccessUuid).Return(nil)
	err := logoutRepo.Logout(dummyTokenDetails[0].AccessToken)
	assert.Nil(suite.T(), err)
}

func (suite *LogoutRepoTestSuite) TestLogout_FailedVerifyAccessToken() {
	logoutRepo := NewLogoutRepository(suite.authMock)
	suite.authMock.On("VerifyAccessToken", dummyTokenDetails[0].AccessToken).Return(authenticator.AccessDetails{}, errors.New("Failed"))
	err := logoutRepo.Logout(dummyTokenDetails[0].AccessToken)
	assert.NotNil(suite.T(), err)
}

func (suite *LogoutRepoTestSuite) TestLogout_FailedDeleteAccessToken() {
	logoutRepo := NewLogoutRepository(suite.authMock)
	suite.authMock.On("VerifyAccessToken", dummyTokenDetails[0].AccessToken).Return(dummyAccessDetails[0], nil)
	suite.authMock.On("DeleteAccessToken", dummyAccessDetails[0].AccessUuid).Return(errors.New("Failed"))
	err := logoutRepo.Logout(dummyTokenDetails[0].AccessToken)
	assert.NotNil(suite.T(), err)
}

func (suite *LogoutRepoTestSuite) SetupTest() {
	suite.authMock = new(authMock)
}

func TestLogoutRepoTestSuite(t *testing.T) {
	suite.Run(t, new(LogoutRepoTestSuite))
}
