package main

import (
	"github.com/benebobaa/amikom-bri-api/app"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/benebobaa/amikom-bri-api/util/mail"
	"github.com/benebobaa/amikom-bri-api/util/token"
	"log"
)

func main() {
	viperConfig, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	tokenMaker := token.NewJWTMaker()
	if err != nil {
		log.Fatalf("Failed to create JWT Maker: %v", err)
	}

	mailSender := mail.NewGmailSender(viperConfig.EmailName, viperConfig.EmailSender, viperConfig.EmailPassword)
	db := app.NewDatabaseConnection(viperConfig.DBDsn)
	validate := app.NewValidator()
	fiber := app.NewFiber(viperConfig)

	app.Bootstrap(&app.BootstrapConfig{
		DB:          db,
		App:         fiber,
		Validate:    validate,
		TokenMaker:  tokenMaker,
		ViperConfig: viperConfig,
		EmailSender: mailSender,
	})

	err = fiber.Listen(":" + viperConfig.PortApp)
	if err != nil {
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
