package main

import (
	"embed"
	"github.com/benebobaa/amikom-bri-api/app"
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/benebobaa/amikom-bri-api/util/mail"
	"github.com/benebobaa/amikom-bri-api/util/token"
	"github.com/gofiber/template/html/v2"
	"log"
	httpLib "net/http"
)

//go:embed resource/*
var resourcefs embed.FS

func main() {
	viperConfig, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	tokenMaker := token.NewJWTMaker()
	if err != nil {
		log.Fatalf("Failed to create JWT Maker: %v", err)
	}

	engine := html.NewFileSystem(httpLib.FS(resourcefs), ".html")

	mailSender := mail.NewTitanSender(viperConfig.EmailName, viperConfig.EmailSender, viperConfig.EmailPassword)
	db := app.NewDatabaseConnection(viperConfig.DBDsn)
	validate := app.NewValidator()
	fiber := app.NewFiber(viperConfig, engine)

	app.Bootstrap(&app.BootstrapConfig{
		DB:          db,
		App:         fiber,
		Validate:    validate,
		TokenMaker:  tokenMaker,
		ViperConfig: viperConfig,
		TitanMail:   mailSender,
	})

	err = fiber.Listen(":" + viperConfig.PortApp)
	if err != nil {
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
