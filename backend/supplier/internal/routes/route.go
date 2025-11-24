package routes

import (
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/handler"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, h handler.SupplierHandler, mw *middleware.AuthMiddleware) {
	api := app.Group("/api/v1/supplier")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "success"})
	})

	// Grup Terproteksi (Harus Login)
	protected := api.Group("/", mw.Protected())

	// READ: Boleh semua role
	protected.Get("/", h.FindAll)
	protected.Get("/:id", h.FindByID)

	// WRITE: Hanya Admin
	adminOnly := protected.Group("/", mw.RequiredRole("admin"))
	adminOnly.Post("/", h.Create)
	adminOnly.Patch("/:id", h.Update)
	adminOnly.Delete("/:id", h.Delete)
}
