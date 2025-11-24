package main

import (
	"log"
	"time"

	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/config"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/database"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/handler"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/middleware"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/queue"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/repository"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/routes"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	log.Println("Memulai supplier-service...")

	// 1. Load Konfigurasi
	cfg := config.LoadConfig()
	log.Println("Konfigurasi di-load.")

	// 2. Inisialisasi Infrastruktur (Database & Broker)
	db := database.InitDatabase(cfg)
	rdb := database.InitRedis(cfg)
	amqpChannel, cleanupRabbitMQ := queue.InitRabbitMQ(cfg)
	defer cleanupRabbitMQ() // Pastikan koneksi ditutup saat main() selesai
	log.Println("Semua database & broker terhubung.")

	// 3. Wiring Layers (Dependency Injection)

	// Repository (Postgres + Redis)
	supplierRepo := repository.NewSupplierRepository(db, rdb)

	// Event Publisher (RabbitMQ)
	eventPublisher := queue.NewRabbitMQPublisher(amqpChannel)

	// Service (Logic + Repo + Publisher)
	supplierSvc := service.NewSupplierService(supplierRepo, eventPublisher)

	// Handler (HTTP -> Service)
	supplierHandler := handler.NewSupplierHandler(supplierSvc)

	// Middleware (Security)
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	log.Println("Semua layer (repo, service, handler, middleware) diinisialisasi.")

	// 4. Inisialisasi Fiber App
	app := fiber.New()

	// 5. Middleware Global
	app.Use(logger.New()) // Logger request

	// Rate Limiter (Mencegah DDOS/Spam)
	app.Use(limiter.New(limiter.Config{
		Max:        10,              // Maks 10 request
		Expiration: 5 * time.Second, // per 5 detik
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Terlalu banyak request, coba lagi nanti.",
			})
		},
	}))

	log.Println("Fiber dan middleware global diinisialisasi.")

	// 6. Setup Routes
	// Kita kirim app, handler, dan middleware keamanan ke router
	routes.SetupRoutes(app, supplierHandler, authMiddleware)
	log.Println("Rute di-setup.")

	// 7. Jalankan Server
	log.Printf("Supplier service siap berjalan di port :%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
