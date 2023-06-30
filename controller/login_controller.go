package controller

import (
	"github.com/febriansr/simple-payment-api/model/app_error"
	entity "github.com/febriansr/simple-payment-api/model/entity"
	"github.com/febriansr/simple-payment-api/usecase"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	BaseController
	router       *gin.RouterGroup
	loginUsecase usecase.LoginUsecase
}

func (l *LoginController) LoginHandler(ctx *gin.Context) {
	var customer entity.Customer

	if err := ctx.ShouldBindJSON(&customer); err != nil {
		l.Failed(ctx, app_error.InvalidError("invalid request body"))
		return
	}

	token, err := l.loginUsecase.Login(customer)

	if err == nil {
		l.Success(ctx, map[string]interface{}{
			"token": token,
		})
	} else {
		l.Failed(ctx, err)
	}
}

func NewLoginController(r *gin.RouterGroup, u usecase.LoginUsecase) *LoginController {
	controller := LoginController{
		loginUsecase: u,
	}
	r.POST("/login", controller.LoginHandler)
	return &controller
}
