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

	//TODO borrar esto
	api.Get("/users/info", func(c *fiber.Ctx) error {
		name, ok := c.Locals("name").(string)
		if !ok {
			return c.Status(400).JSON(fiber.Map{
				"message": "Error getting email",
			})
		}
		return c.SendString("Hello, World! " + name)
	})

	api.Get("/users/:_id", func(c *fiber.Ctx) error {
		return c.SendStatus(201)
	})
}
