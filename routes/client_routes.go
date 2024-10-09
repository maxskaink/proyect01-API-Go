package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxskaink/proyect01-api-go/controllers"
	"github.com/maxskaink/proyect01-api-go/middlewares/auth"
	"github.com/maxskaink/proyect01-api-go/services"
)

// UserRoutes asign the routes for the users
// each route have a controller
func UserRoutes(app *fiber.App, serviceU *services.UsersService) {
	controllers := controllers.NewClientControllers(serviceU)

	api := app.Group("/api")
	api.Get("/users", controllers.GetAllUsers)
	api.Post("/users", controllers.CreateUser)
	api.Post("/users/login", controllers.LogInUser)

	api.Use(auth.UserAuth())

	api.Get("/users/:id", controllers.GetUserbyID)
	api.Put("/users/:id", controllers.UpdateUser)
	api.Patch("/users/:id", controllers.PatchUser)
	api.Delete("/users/:id", controllers.DeleteUser)
}
