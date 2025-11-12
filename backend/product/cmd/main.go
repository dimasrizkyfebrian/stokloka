package main

import (
	"log"
	"time"

	"github.com/dimasrizkyfebrian/stokloka/product/internal/config"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/database"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/handler"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/middleware"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/queue"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/repository"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/routes"
	"github.com/dimasrizkyfebrian/stokloka/product/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	log.Println("Memulai product-service...")

	// Load Konfigurasi
	cfg := config.LoadConfig()
	log.Println("Konfigurasi di-load.")

	// Inisialisasi Database (Postgres, Redis, RabbitMQ)
	db := database.InitDatabase(cfg)
	rdb := database.InitRedis(cfg)
	amqpChannel, cleanupRabbitMQ := queue.InitRabbitMQ(cfg)
	defer cleanupRabbitMQ()
	log.Println("Semua database & broker terhubung.")

	// Inisialisasi Layers
	productRepo := repository.NewProductRepository(db, rdb)
	eventPublisher := queue.NewRabbitMQPublisher(amqpChannel)
	productSvc := service.NewProductService(productRepo, eventPublisher)
	productHandler := handler.NewProductHandler(productSvc)

	// Inisialisasi Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg)
	log.Println("Semua layer (repo, service, handler, middleware) diinisialisasi.")

	// Inisialisasi Fiber App
	app := fiber.New()
	app.Use(limiter.New(limiter.Config{
		Max:        10,              // Izinkan 10 request
		Expiration: 5 * time.Second, // per 5 detik
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Terlalu banyak request, coba lagi nanti.",
			})
		},
	}))

	log.Println("Fiber dan middleware global diinisialisasi.")

	// Setup Rute
	routes.SetupRoutes(app, productHandler, authMiddleware)
	log.Println("Rute di-setup.")

	// Jalankan server
	log.Printf("Product service siap berjalan di port :%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
