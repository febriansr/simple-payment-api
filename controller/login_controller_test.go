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

var dummyCustomer = []entity.Customer{
	{
		Username: "dummyUsername",
		Password: "dummyPassword",
	},
}

var dummyTokenDetails = []authenticator.TokenDetails{
	{
		AccessToken: "Dummy Access Token",
		AccessUuid:  "Dummy Access Uuid",
		AtExpires:   5,
	},
}

type LoginUsecaseMock struct {
	mock.Mock
}

func (l *LoginUsecaseMock) Login(customer entity.Customer) (token string, err error) {
	args := l.Called(customer)
	if args.Get(0) == nil {
		return "", errors.New("Failed")
	}
	return args.Get(0).(string), args.Error(1)
}

type LoginControllerTestSuite struct {
	suite.Suite
	routerMock      *gin.Engine
	routerGroupMock *gin.RouterGroup
	usecaseMock     *LoginUsecaseMock
}

func (suite *LoginControllerTestSuite) TestLogin_Success() {
	customer := dummyCustomer[0]
	NewLoginController(suite.routerGroupMock, suite.usecaseMock)
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(customer)
	request, _ := http.NewRequest(http.MethodPost, "/v1/login", bytes.NewBuffer(reqBody))

	suite.usecaseMock.On("Login", customer).Return(dummyTokenDetails[0].AccessToken, nil)
	suite.routerMock.ServeHTTP(r, request)
	var response res.ApiResponse
	json.Unmarshal(r.Body.Bytes(), &response)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	assert.Equal(suite.T(), response.Data, map[string]interface{}{
		"token": dummyTokenDetails[0].AccessToken,
	})
}

func (suite *LoginControllerTestSuite) TestLogin_FailedBindJSON() {
	invalidReqBody := []byte(`{1}`)
	NewLoginController(suite.routerGroupMock, suite.usecaseMock)
	r := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/v1/login", bytes.NewBuffer(invalidReqBody))

	suite.routerMock.ServeHTTP(r, request)
	var response res.ApiResponse
	json.Unmarshal(r.Body.Bytes(), &response)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *LoginControllerTestSuite) TestLogin_FailedErrorUsecase() {
	customer := dummyCustomer[0]
	NewLoginController(suite.routerGroupMock, suite.usecaseMock)
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(customer)
	request, _ := http.NewRequest(http.MethodPost, "/v1/login", bytes.NewBuffer(reqBody))

	suite.usecaseMock.On("Login", customer).Return("", errors.New("Failed"))
	suite.routerMock.ServeHTTP(r, request)
	var response res.ApiResponse
	json.Unmarshal(r.Body.Bytes(), &response)

	assert.Equal(suite.T(), "XX", response.Status.Code)
}

func (suite *LoginControllerTestSuite) SetupTest() {
	suite.routerMock = gin.Default()
	suite.routerGroupMock = suite.routerMock.Group("/v1")
	suite.usecaseMock = new(LoginUsecaseMock)
}

func TestLoginControllerTestSuite(t *testing.T) {
	suite.Run(t, new(LoginControllerTestSuite))
}
