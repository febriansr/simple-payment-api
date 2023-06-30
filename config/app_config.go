package config

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/febriansr/simple-payment-api/utils"
)

type ApiConfig struct {
	ServerPort string
	ServerHost string
}

type JsonFileConfig struct {
	Customer string
	Merchant string
	History  string
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     string
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifetime time.Duration
}

type RedisConfig struct {
	Address  string
	Password string
	Db       int
}

type AppConfig struct {
	ApiConfig
	JsonFileConfig
	TokenConfig
	RedisConfig
}

func (c *AppConfig) readConfigFile() {
	envFilePath := ".env"
	c.JsonFileConfig = JsonFileConfig{
		Customer: utils.DotEnv("JSON_FILE_NAME_CUSTOMER", envFilePath),
		Merchant: utils.DotEnv("JSON_FILE_NAME_MERCHANT", envFilePath),
		History:  utils.DotEnv("JSON_FILE_NAME_HISTORY", envFilePath),
	}
	c.ApiConfig = ApiConfig{
		ServerPort: utils.DotEnv("SERVER_PORT", envFilePath),
		ServerHost: utils.DotEnv("SERVER_HOST", envFilePath),
	}
	lifeTime, _ := strconv.Atoi(utils.DotEnv("ACCESS_TOKEN_LIFETIME", envFilePath))
	c.TokenConfig = TokenConfig{
		ApplicationName:     utils.DotEnv("APPLICATION_NAME", envFilePath),
		JwtSignatureKey:     utils.DotEnv("JWT_SIGNATURE_KEY", envFilePath),
		JwtSigningMethod:    jwt.SigningMethodHS256,
		AccessTokenLifetime: time.Duration(lifeTime) * time.Minute,
	}
	c.RedisConfig = RedisConfig{
		Address:  utils.DotEnv("REDDIS_ADDRESS", envFilePath),
		Password: utils.DotEnv("REDDIS_PASSWORD", envFilePath),
		Db:       0,
	}
}

func NewConfig() AppConfig {
	config := AppConfig{}
	config.readConfigFile()
	return config
}
