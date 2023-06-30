package controller

import (
	"github.com/febriansr/simple-payment-api/usecase"
	"github.com/febriansr/simple-payment-api/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type LogoutController struct {
	logoutUsecase usecase.LogoutUsecase
	BaseController
	router *gin.RouterGroup
}

func (c *LogoutController) LogoutHandler(ctx *gin.Context) {
	token, err := authenticator.BindAuthHeader(ctx)
	if err != nil {
		c.Failed(ctx, err)
		ctx.Abort()
		return
	}

	err = c.logoutUsecase.Logout(token)

	if err != nil {
		c.Failed(ctx, err)
		return
	}
	c.Success(ctx, nil)
}

func NewLogoutController(r *gin.RouterGroup, u usecase.LogoutUsecase) *LogoutController {
	controller := LogoutController{
		logoutUsecase: u,
	}
	r.POST("/logout", controller.LogoutHandler)
	return &controller
}
