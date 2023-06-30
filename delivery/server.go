package delivery

import (
	"fmt"
	"log"

	"github.com/febriansr/simple-payment-api/config"
	"github.com/febriansr/simple-payment-api/controller"
	"github.com/febriansr/simple-payment-api/manager"
	"github.com/febriansr/simple-payment-api/middleware"
	"github.com/febriansr/simple-payment-api/utils/authenticator"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type AppServer struct {
	usecaseManager manager.UsecaseManager
	authenticator  authenticator.AccessToken
	engine         *gin.Engine
	host           string
}

func (p *AppServer) menu() {
	routes := p.engine.Group("/v1")
	routes.Use(middleware.LoggingMiddleware(".log"))
	middleware := middleware.NewAuthTokenMiddleware(p.authenticator)
	p.loginController(routes)
	p.logoutController(routes)
	p.paymentController(routes, p.authenticator, middleware)
}

func (p *AppServer) loginController(rg *gin.RouterGroup) {
	controller.NewLoginController(rg, p.usecaseManager.LoginUsecase())
}

func (p *AppServer) logoutController(rg *gin.RouterGroup) {
	controller.NewLogoutController(rg, p.usecaseManager.LogoutUsecase())
}

func (p *AppServer) paymentController(rg *gin.RouterGroup, authenticator authenticator.AccessToken, middleware middleware.AuthTokenMiddleware) {
	controller.NewPaymentController(rg, p.usecaseManager.PaymentUsecase(), authenticator, middleware)
}

func (p *AppServer) Run() {
	p.menu()
	err := p.engine.Run(p.host)
	defer func() {
		if err := recover(); err != nil {
			log.Println("Application failed to run", err)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
}

func Server() *AppServer {
	router := gin.Default()
	config := config.NewConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Address,
		Password: config.RedisConfig.Password,
		DB:       config.RedisConfig.Db,
	})
	authenticator := authenticator.NewAccessToken(config.TokenConfig, client)
	repositoryManager := manager.NewRepositoryManager(config.JsonFileConfig, authenticator)
	usecaseManager := manager.NewUsecaseManager(repositoryManager, authenticator)
	host := fmt.Sprintf("%s:%s", config.ServerHost, config.ServerPort)
	return &AppServer{
		usecaseManager: usecaseManager,
		engine:         router,
		host:           host,
		authenticator:  authenticator,
	}
}
