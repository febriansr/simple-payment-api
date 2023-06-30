package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/febriansr/simple-payment-api/model/dto/res"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LogoutUsecaseMock struct {
	mock.Mock
}

func (l *LogoutUsecaseMock) Logout(token string) error {
	args := l.Called(token)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

type BindAuthHeaderMock struct {
	mock.Mock
}

func (b *BindAuthHeaderMock) BindAuthHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Empty token")
	}
	return dummyTokenDetails[0].AccessToken, nil
}

type LogoutControllerTestSuite struct {
	suite.Suite
	routerMock         *gin.Engine
	routerGroupMock    *gin.RouterGroup
	usecaseMock        *LogoutUsecaseMock
	bindAuthHeaderMock *BindAuthHeaderMock
}

func (suite *LogoutControllerTestSuite) TestLogout_Success() {
	logoutController := NewLogoutController(suite.routerGroupMock, suite.usecaseMock)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/v1/logout", bytes.NewBuffer([]byte(`{}`)))
	request.Header.Set("Authorization", dummyTokenDetails[0].AccessToken)
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = request
	suite.bindAuthHeaderMock.On("BindAuthHeader", ctx).Return(dummyTokenDetails[0].AccessToken, nil)
	suite.usecaseMock.On("Logout", dummyTokenDetails[0].AccessToken).Return(nil)
	logoutController.LogoutHandler(ctx)

	var response res.ApiResponse
	json.Unmarshal(r.Body.Bytes(), &response)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
}

func (suite *LogoutControllerTestSuite) TestLogout_FailedUsecase() {
	logoutController := NewLogoutController(suite.routerGroupMock, suite.usecaseMock)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/v1/logout", bytes.NewBuffer([]byte(`{}`)))
	request.Header.Set("Authorization", dummyTokenDetails[0].AccessToken)
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = request
	suite.bindAuthHeaderMock.On("BindAuthHeader", ctx).Return(dummyTokenDetails[0].AccessToken, nil)
	suite.usecaseMock.On("Logout", dummyTokenDetails[0].AccessToken).Return(errors.New("Failed"))
	logoutController.LogoutHandler(ctx)

	var response res.ApiResponse
	json.Unmarshal(r.Body.Bytes(), &response)

	assert.Equal(suite.T(), "XX", response.Status.Code)
}

func (suite *LogoutControllerTestSuite) TestLogout_FailedBindHeader() {
	logoutController := NewLogoutController(suite.routerGroupMock, suite.usecaseMock)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/v1/logout", bytes.NewBuffer([]byte(`{}`)))
	request.Header.Set("Authorization", "Bearer")
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = request
	suite.bindAuthHeaderMock.On("BindAuthHeader", ctx).Return("", errors.New(""))
	logoutController.LogoutHandler(ctx)

	var response res.ApiResponse
	json.Unmarshal(r.Body.Bytes(), &response)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *LogoutControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.routerGroupMock = suite.routerMock.Group("/v1")
	suite.usecaseMock = new(LogoutUsecaseMock)
	suite.bindAuthHeaderMock = new(BindAuthHeaderMock)
}

func TestLogoutControllerTestSuite(t *testing.T) {
	suite.Run(t, new(LogoutControllerTestSuite))
}
