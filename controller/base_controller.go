package controller

import (
	"github.com/febriansr/simple-payment-api/model/dto/res"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

func (b *BaseController) Success(ctx *gin.Context, data any) {
	res.NewSuccessJsonResponse(ctx, data).Send()
}

func (b *BaseController) Failed(ctx *gin.Context, err error) {
	res.NewErrorJsonResponse(ctx, err).Send()
}
