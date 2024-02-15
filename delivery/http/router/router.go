package router

import (
	"github.com/benebobaa/amikom-bri-api/delivery/http/controller"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                    *fiber.App
	AuthMiddleware         fiber.Handler
	UserController         controller.UserController
	WebController          controller.WebController
	TransferController     controller.TransferController
	EntryController        controller.EntryController
	ExpensesController     controller.ExpensesPlanController
	NotificationController controller.NotificationController
	AuthController         controller.AuthController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {

	// Register
	c.App.Post("/api/v1/auth/_register", c.AuthController.Register)

	// Login
	c.App.Post("/api/v1/auth/_login", c.AuthController.Login)

	// Verify Email
	c.App.Get("/api/v1/auth/_verify-email", c.WebController.VerifyEmail)

	// Forgot Password
	c.App.Post("/api/v1/auth/_forgot-password", c.UserController.ForgotPassword)

	c.App.Get("/users/reset-password", c.WebController.ResetPasswordForm)
	c.App.Post("/users/reset-password-submit", c.WebController.ResetPasswordSubmit)

	c.App.Get("/users/reset-password-success", c.WebController.ResetPasswordSuccess)

	// Renew Access Token
	c.App.Post("/api/v1/auth/_renew-token", c.AuthController.RenewAccessToken)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)

	// User
	c.App.Delete("/api/v1/users/:username", c.UserController.Delete)
	c.App.Get("/api/v1/users/profile", c.UserController.Profile)
	c.App.Get("/api/v1/users", c.UserController.FindAll)
	c.App.Patch("/api/v1/users", c.UserController.Update)

	// Transfer
	c.App.Post("/api/v1/transfer", c.TransferController.Transfer)

	// Entries
	c.App.Get("/api/v1/entries", c.EntryController.FindAllEntries)
	c.App.Get("/api/v1/entries/filter", c.EntryController.FindAllFilter)

	// Expenses Plan
	c.App.Post("/api/v1/expenses-plan", c.ExpensesController.Create)
	c.App.Delete("/api/v1/expenses-plan/:id", c.ExpensesController.Delete)
	c.App.Get("/api/v1/expenses-plan", c.ExpensesController.FindAllFilter)
	c.App.Patch("/api/v1/expenses-plan/:id", c.ExpensesController.Update)
	c.App.Get("/api/v1/expenses-plan/export-pdf", c.ExpensesController.FindAllExportPdf)

	// Notification
	c.App.Get("/api/v1/notifications", c.NotificationController.FindALL)
}
