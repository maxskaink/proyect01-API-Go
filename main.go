package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/maxskaink/proyect01-api-go/routes"
)

func main() {
	app := fiber.New()

	// Routes
	routes.APIRoutes(app)

	// Init API
	loadENV()
	PORT_API := os.Getenv("PORT_API")
	fmt.Println("Server running on port: " + PORT_API)
	log.Fatal(app.Listen(":" + PORT_API))
}

func loadENV() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
