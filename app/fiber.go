package app

import (
	"github.com/benebobaa/amikom-bri-api/util"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
)

func NewFiber(c util.Config, engine *html.Engine) *fiber.App {
	config := fiber.Config{
		AppName:           c.AppName,
		EnablePrintRoutes: true,
		JSONEncoder:       sonic.Marshal,
		JSONDecoder:       sonic.Unmarshal,
		Views:             engine,
	}

	if c.GOEnv == "production" {
		config.Prefork = true
	}

	app := fiber.New(config)
	app.Use(logger.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	return app
}
