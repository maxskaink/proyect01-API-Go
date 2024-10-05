package main

import (
	"context"
	"log"
	"os"

	"github.com/maxskaink/proyect01-api-go/services"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/maxskaink/proyect01-api-go/routes"
)

// main entry point for mi API rest aplication
// get de eviroment variables, connect with de database
// get de routes for de api and init the litening of the API
func main() {
	loadENV()
	client := services.InitDataBase()
	defer client.Disconnect(context.Background())
	app := fiber.New()

	// Routes
	routes.APIRoutes(app)
	routes.UserRoutes(app)

	// Init API
	PORT_API := os.Getenv("PORT_API")

	log.Fatal(app.Listen(":" + PORT_API))
}

// loadENV get to the so variable, all the enviroments variables
// of the .env file
func loadENV() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
