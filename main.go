package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()
	loadENV()
	PORT_API := os.Getenv("PORT_API")

	log.Fatal(app.Listen(":" + PORT_API))
}

func loadENV() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
