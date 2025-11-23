package main

import (
	"log"

	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/config"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/database"
	"github.com/dimasrizkyfebrian/stokloka/supplier/internal/queue"
	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("Memulai supplier-service...")

	cfg := config.LoadConfig()

	// Inisialisasi Infrastruktur
	_ = database.InitDatabase(cfg)        // Postgres & Migrasi
	_ = database.InitRedis(cfg)           // Redis
	_, cleanup := queue.InitRabbitMQ(cfg) // RabbitMQ
	defer cleanup()

	app := fiber.New()

	app.Get("/api/v1/supplier/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "success", "message": "Supplier Service Running!"})
	})

	log.Printf("Supplier service berjalan di port :%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
