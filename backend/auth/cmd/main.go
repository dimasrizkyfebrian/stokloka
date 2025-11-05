package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/dimasrizkyfebrian/stokloka/auth/internal/routes"
)

func main() {
	app := fiber.New()

	routes.SetupRoutes(app)

	log.Println("Starting Auth service on port 8080")
	err := app.Listen(":8080")
	if err != nil {
		log.Fatalf("Error starting Auth service: %v", err)
	}
}
