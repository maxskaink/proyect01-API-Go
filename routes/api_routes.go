package routes

import (
	"github.com/maxskaink/proyect01-api-go/controllers"

	"github.com/gofiber/fiber/v2"
)

func APIRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/ping", controllers.PingAPI)
}
