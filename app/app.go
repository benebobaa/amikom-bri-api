package app

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/controller"
	"github.com/benebobaa/amikom-bri-api/delivery/http/middleware"
	"github.com/benebobaa/amikom-bri-api/delivery/http/router"
	"github.com/benebobaa/amikom-bri-api/domain/repository"
	"github.com/benebobaa/amikom-bri-api/domain/usecase"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/benebobaa/amikom-bri-api/util/mail"
	"github.com/benebobaa/amikom-bri-api/util/token"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB          *gorm.DB
	App         *fiber.App
	Validate    *validator.Validate
	TokenMaker  token.Maker
	ViperConfig util.Config
	EmailSender mail.EmailSender
}

func Bootstrap(config *BootstrapConfig) {

	// Setup repositories
	userRepository := repository.NewUserRepository()
	emailRepository := repository.NewEmailRepository()
	sessionRepository := repository.NewSessionRepository()
	accountRepository := repository.NewAccountRepository()

	// Setup usecases
	userUsecase := usecase.NewUserUsecase(config.DB, config.Validate, userRepository, emailRepository, accountRepository)
	loginUsecase := usecase.NewLoginUseCase(config.DB, config.Validate, config.TokenMaker, config.ViperConfig, userRepository, sessionRepository)

	// Setup controller
	userController := controller.NewUserController(userUsecase, loginUsecase)

	// Setup middleware
	authMiddleware := middleware.AuthMiddleware(config.TokenMaker, config.ViperConfig)

	routeConfig := router.RouteConfig{
		App:            config.App,
		AuthMiddleware: authMiddleware,
		UserController: userController,
	}

	routeConfig.Setup()
}