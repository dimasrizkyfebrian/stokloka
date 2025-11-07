package routes

import (
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/handler"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes mendaftarkan rute untuk auth
func SetupAuthRoutes(app *fiber.App, authHandler handler.AuthHandler, mw *middleware.AuthMiddleware) {
	// Grup /api/v1/auth
	api := app.Group("/api/v1/auth")

	// Endpoint publik
	api.Post("/login", authHandler.Login)

	// --- Endpoint Terproteksi ---
	protected := api.Group("/", mw.Protected())

	protected.Post(
		"/register",
		mw.RequiredRole("admin"), // Cek apakah role == "admin"
		authHandler.Register,     // Jika iya maka baru jalankan handler register
	)

	// Health check (sesuai blueprint)
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Auth service is running!",
		})
	})
}
