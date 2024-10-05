package routes

import (
	"github.com/maxskaink/proyect01-api-go/controllers"

	"github.com/gofiber/fiber/v2"
)

// APIRoutes asign the api routes to an fiber app
// each routes have a controller for the endpoint
func APIRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/ping", controllers.PingAPI)
}
