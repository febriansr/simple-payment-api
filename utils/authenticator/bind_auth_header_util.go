package authenticator

import (
	"strings"

	"github.com/febriansr/simple-payment-api/model/app_error"
	"github.com/febriansr/simple-payment-api/model/dto/req"
	"github.com/gin-gonic/gin"
)

func BindAuthHeader(c *gin.Context) (string, error) {
	header := new(req.AuthHeader)
	if err := c.ShouldBindHeader(header); err != nil {
		return "", app_error.InternalServerError(err.Error())
	}
	tokenString := strings.Replace(header.AuthorizationHeader, "Bearer ", "", -1)

	if tokenString == "Bearer" || tokenString == "" {
		return "", app_error.InvalidError("Empty token")
	}
	return tokenString, nil
}
