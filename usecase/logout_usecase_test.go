package usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type logoutRepoMock struct {
	mock.Mock
}

func (l *logoutRepoMock) Logout(token string) error {
	args := l.Called(token)
	if args[0] != nil {
		return errors.New("Failed")
	}
	return nil
}

type LogoutUsecaseTestSuite struct {
	logoutRepoMock *logoutRepoMock
	suite.Suite
}

func (suite *LogoutUsecaseTestSuite) TestLogout_Success() {
	logoutUsecase := NewLogoutUsecase(suite.logoutRepoMock)
	suite.logoutRepoMock.On("Logout", dummyTokenDetails[0].AccessToken).Return(nil)
	err := logoutUsecase.Logout(dummyTokenDetails[0].AccessToken)
	assert.Nil(suite.T(), err)
}

func (suite *LogoutUsecaseTestSuite) TestLogout_Failed() {
	logoutUsecase := NewLogoutUsecase(suite.logoutRepoMock)
	suite.logoutRepoMock.On("Logout", dummyTokenDetails[0].AccessToken).Return(errors.New("Failed"))
	err := logoutUsecase.Logout(dummyTokenDetails[0].AccessToken)
	assert.NotNil(suite.T(), err)
}

func (suite *LogoutUsecaseTestSuite) SetupTest() {
	suite.logoutRepoMock = new(logoutRepoMock)
}

func TestLogoutUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(LogoutUsecaseTestSuite))
}
