package controller

import (
	"github.com/febriansr/simple-payment-api/middleware"
	"github.com/febriansr/simple-payment-api/model/app_error"
	entity "github.com/febriansr/simple-payment-api/model/entity"
	"github.com/febriansr/simple-payment-api/usecase"
	"github.com/febriansr/simple-payment-api/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentUsecase usecase.PaymentUsecase
	authenticator  authenticator.AccessToken
	BaseController
	router *gin.RouterGroup
}

func (l *PaymentController) PaymentHandler(ctx *gin.Context) {
	var transaction entity.History

	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		l.Failed(ctx, app_error.InvalidError(err.Error()))
		return
	}

	token, err := authenticator.BindAuthHeader(ctx)
	if err != nil {
		l.Failed(ctx, err)
		ctx.Abort()
		return
	}

	accountDetails, err := l.authenticator.VerifyAccessToken(token)
	if err != nil {
		l.Failed(ctx, err)
		return
	}

	transaction.CustomerUsername = accountDetails.Username

	err = l.paymentUsecase.PayTransaction(transaction)

	if err == nil {
		l.Success(ctx, nil)
	} else {
		l.Failed(ctx, err)
	}
}

func NewPaymentController(r *gin.RouterGroup, u usecase.PaymentUsecase, a authenticator.AccessToken, m middleware.AuthTokenMiddleware) *PaymentController {
	controller := PaymentController{
		paymentUsecase: u,
		authenticator:  a,
	}
	rm := r.Group("/menu", m.RequireToken())
	rm.POST("/payment", controller.PaymentHandler)
	return &controller
}
