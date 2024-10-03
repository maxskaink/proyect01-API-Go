package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxskaink/proyect01-api-go/controllers"
)

func UserRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/users", controllers.GetAllUsers)
	api.Post("/users", controllers.CreateUser)

}
