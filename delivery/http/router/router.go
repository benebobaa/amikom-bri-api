package router

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/controller"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	AuthMiddleware     fiber.Handler
	UserController     controller.UserController
	WebController      controller.WebController
	TransferController controller.TransferController
	EntryController    controller.EntryController
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
	c.App.Get("/api/v1/auth/_verify-email", c.WebController.VerifyEmail)

	// Forgot Password
	c.App.Post("/api/v1/auth/_forgot-password", c.UserController.ForgotPassword)

	c.App.Get("/users/reset-password", c.WebController.ResetPasswordForm)
	c.App.Post("/users/reset-password-submit", c.WebController.ResetPasswordSubmit)

	c.App.Get("/users/reset-password-success", c.WebController.ResetPasswordSuccess)
	//c.App.Get("/_test", c.WebController.Test)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)

	// User
	c.App.Delete("/api/v1/users/:username", c.UserController.Delete)
	c.App.Get("/api/v1/users/profile", c.UserController.Profile)
	c.App.Get("/api/v1/users", c.UserController.FindAll)

	// Transfer
	c.App.Post("/api/v1/transfer", c.TransferController.Transfer)

	// Entries
	c.App.Get("/api/v1/entries", c.EntryController.FindAllEntries)
}
