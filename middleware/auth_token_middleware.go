package middleware

import (
	"github.com/febriansr/simple-payment-api/model/dto/res"
	"github.com/febriansr/simple-payment-api/utils/authenticator"
	"github.com/gin-gonic/gin"
)

type AuthTokenMiddleware interface {
	RequireToken() gin.HandlerFunc
}

type authTokenMiddleware struct {
	authenticator authenticator.AccessToken
}

func (a *authTokenMiddleware) RequireToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := authenticator.BindAuthHeader(ctx)
		if err != nil {
			res.NewErrorJsonResponse(ctx, err).Send()
			ctx.Abort()
			return
		}
		accountDetails, err := a.authenticator.VerifyAccessToken(token)
		if err != nil {
			res.NewErrorJsonResponse(ctx, err).Send()
			ctx.Abort()
			return
		}

		err = a.authenticator.FetchAccessToken(accountDetails)
		if err != nil {
			res.NewErrorJsonResponse(ctx, err).Send()
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func NewAuthTokenMiddleware(authenticator authenticator.AccessToken) AuthTokenMiddleware {
	return &authTokenMiddleware{
		authenticator: authenticator,
	}
}
