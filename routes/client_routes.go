package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxskaink/proyect01-api-go/controllers"
	"github.com/maxskaink/proyect01-api-go/middlewares/auth"
)

func UserRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/users", controllers.GetAllUsers)
	api.Post("/users", controllers.CreateUser)
	api.Post("/users/login", controllers.LogInUser)

	api.Use(auth.UserAuth())

	api.Get("/users/:id", controllers.GetUserbyID)
	api.Put("/users/:id", controllers.UpdateUser)
}
