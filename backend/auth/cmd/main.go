package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/api/v1/auth/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Auth service is up and running!",
		})
	})

	log.Println("Starting Auth service on port 8080")
	err := app.Listen(":8080")
	if err != nil {
		log.Fatalf("Error starting Auth service: %v", err)
	}
}
