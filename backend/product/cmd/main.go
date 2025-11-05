package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/dimasrizkyfebrian/stokloka/product/internal/routes"
)

func main() {
	app := fiber.New()

	routes.SetupRoutes(app)

	log.Println("Starting Product service on port 8081")
	err := app.Listen(":8081")
	if err != nil {
		log.Fatalf("Error starting Product service: %v", err)
	}
}
