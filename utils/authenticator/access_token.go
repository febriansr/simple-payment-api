package authenticator

import (
	"context"
	"time"

	"github.com/febriansr/simple-payment-api/model/app_error"

	"github.com/dgrijalva/jwt-go"
	"github.com/febriansr/simple-payment-api/config"
	entity "github.com/febriansr/simple-payment-api/model/entity"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type AccessToken interface {
	CreateAccessToken(customer *entity.Customer) (TokenDetails, error)
	VerifyAccessToken(tokenString string) (AccessDetails, error)
	StoreAccessToken(username string, tokenDetails TokenDetails) error
	FetchAccessToken(accessDetails AccessDetails) error
	DeleteAccessToken(accessUuid string) error
}

type accessToken struct {
	config config.TokenConfig
	client *redis.Client
}

func (t *accessToken) CreateAccessToken(customer *entity.Customer) (TokenDetails, error) {
	tokenDetails := TokenDetails{}
	now := time.Now().UTC()
	end := now.Add(t.config.AccessTokenLifetime)
	tokenDetails.AtExpires = end.Unix()
	tokenDetails.AccessUuid = uuid.New().String()
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer: t.config.ApplicationName,
		},
		Username:   customer.Username,
		AccessUuid: tokenDetails.AccessUuid,
	}
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = end.Unix()
	token := jwt.NewWithClaims(
		t.config.JwtSigningMethod,
		claims,
	)
	newToken, err := token.SignedString([]byte(t.config.JwtSignatureKey))
	if err != nil {
		return tokenDetails, app_error.InternalServerError("Failed to sign access token: " + err.Error())
	}
	tokenDetails.AccessToken = newToken
	return tokenDetails, nil
}

func (t *accessToken) VerifyAccessToken(tokenString string) (AccessDetails, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, app_error.Unauthorized("Invalid signing method")
		} else if method != t.config.JwtSigningMethod {
			return nil, app_error.Unauthorized("Invalid signing method")
		}
		return []byte(t.config.JwtSignatureKey), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	accessDetails := AccessDetails{}
	if !ok || !token.Valid || claims["iss"] != t.config.ApplicationName {
		return accessDetails, app_error.Unauthorized("Invalid access method")
	}
	username := claims["username"].(string)
	accessDetails.AccessUuid = claims["AccessUuid"].(string)
	accessDetails.Username = username
	return accessDetails, nil
}

func (t *accessToken) StoreAccessToken(username string, tokenDetails TokenDetails) error {
	end := time.Unix(tokenDetails.AtExpires, 0)
	now := time.Now()
	err := t.client.Set(context.Background(), tokenDetails.AccessUuid, username, end.Sub(now)).Err()
	if err != nil {
		return app_error.InternalServerError("Failed to store access token: " + err.Error())
	}
	return nil
}

func (t *accessToken) FetchAccessToken(accessDetails AccessDetails) error {
	username, err := t.client.Get(context.Background(), accessDetails.AccessUuid).Result()
	if err != nil {
		return app_error.Unauthorized("Failed to fetch access token: " + err.Error())
	}
	if username == "" {
		return app_error.Unauthorized("Invalid token")
	}
	return nil
}

func (t *accessToken) DeleteAccessToken(accessUuid string) error {
	rowsAffected, err := t.client.Del(context.Background(), accessUuid).Result()
	if err != nil {
		return app_error.InternalServerError("Failed to delete access token: " + err.Error())
	}
	if rowsAffected == 0 {
		return app_error.DataNotFound("Invalid token")
	}
	return nil
}

func NewAccessToken(config config.TokenConfig, client *redis.Client) AccessToken {
	return &accessToken{
		config: config,
		client: client,
	}
}
