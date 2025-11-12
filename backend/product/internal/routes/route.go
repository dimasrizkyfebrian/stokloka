package routes

import (
	"github.com/dimasrizkyfebrian/stokloka/product/internal/handler"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, productHandler handler.ProductHandler, mw *middleware.AuthMiddleware) {

	api := app.Group("/api/v1/product")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success", "message": "Product service is up and running!",
		})
	})

	// --- Grup Terproteksi (Semua butuh login) ---
	// Semua rute di bawah ini butuh token JWT yang valid (Admin, Manager, Staff)
	protected := api.Group("/", mw.Protected())

	// Rute READ (GET) - Boleh diakses semua role yang login
	protected.Get("/categories", productHandler.GetAllCategories)
	protected.Get("/units", productHandler.GetAllUnits)
	protected.Get("/products", productHandler.GetAllProducts)
	protected.Get("/products/:id", productHandler.GetProductByID)

	// --- Grup Admin (Hanya Admin) ---
	// Rute CRUD (POST, PATCH, DELETE) - Hanya boleh 'admin'
	adminOnly := protected.Group("/", mw.RequiredRole("admin"))

	// Kategori
	adminOnly.Post("/categories", productHandler.CreateCategory)
	adminOnly.Patch("/categories/:id", productHandler.UpdateCategory)
	adminOnly.Delete("/categories/:id", productHandler.DeleteCategory)

	// Unit
	adminOnly.Post("/units", productHandler.CreateUnit)
	adminOnly.Patch("/units/:id", productHandler.UpdateUnit)
	adminOnly.Delete("/units/:id", productHandler.DeleteUnit)

	// Produk
	adminOnly.Post("/products", productHandler.CreateProduct)
	adminOnly.Patch("/products/:id", productHandler.UpdateProduct)
	adminOnly.Delete("/products/:id", productHandler.DeleteProduct)
}
