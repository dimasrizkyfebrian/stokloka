package main

import (
	"log"
	"time"

	"github.com/dimasrizkyfebrian/stokloka/auth/internal/config"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/database"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/handler"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/middleware"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/repository"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/routes"
	"github.com/dimasrizkyfebrian/stokloka/auth/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	log.Println("Memulai auth-service...")

	// 1. Load Konfigurasi
	cfg := config.LoadConfig()
	log.Println("Konfigurasi di-load.")

	// 2. Inisialisasi Database (Konek, Migrasi, Seed)
	db := database.InitDatabase(cfg)
	log.Println("Database terhubung.")

	// 3. Inisialisasi Layers (Dependency Injection)
	authRepo := repository.NewAuthRepository(db)
	// Kirim JWTSecret dari config ke service
	authSvc := service.NewAuthService(authRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authSvc)
	authMiddleware := middleware.NewAuthMiddleware(cfg)
	log.Println("Semua layer (repo, service, handler, middleware) diinisialisasi.")

	// 4. Inisialisasi Fiber App
	app := fiber.New()
	app.Use(logger.New()) // Tambahkan logger middleware
	log.Println("Fiber diinisialisasi.")

	// 5. Tambahkan Global Middleware
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{
		Max:        5,                // 5 request
		Expiration: 10 * time.Second, // per 10 detik
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // gunakan IP sebagai key
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Terlalu banyak request, coba lagi nanti",
			})
		},
	}))

	// 6. Setup Rute
	routes.SetupAuthRoutes(app, authHandler, authMiddleware)
	log.Println("Rute di-setup.")

	// 7. Jalankan server
	log.Printf("Auth service siap berjalan di port :%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
