package routes

import (
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/handler"
	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes mendaftarkan rute untuk auth
func SetupAuthRoutes(app *fiber.App, authHandler handler.AuthHandler) {
	// Grup /api/v1/auth
	api := app.Group("/api/v1/auth")

	// Endpoint publik
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// Health check (sesuai blueprint)
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Auth service is running!",
		})
	})
}
