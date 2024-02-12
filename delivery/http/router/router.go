package router

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/controller"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler
	UserController controller.UserController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {

	// Register
	c.App.Post("/api/v1/auth/_register", c.UserController.Register)

	// Login
	c.App.Post("/api/v1/auth/_login", c.UserController.Login)

	// Verify Email
	c.App.Get("/api/v1/auth/_verify-email", c.UserController.VerifyEmail)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)

}
