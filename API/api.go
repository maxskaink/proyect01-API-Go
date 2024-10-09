package api

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/maxskaink/proyect01-api-go/routes"
	"github.com/maxskaink/proyect01-api-go/services"
)

type API struct {
	app  *fiber.App
	PORT string
}

func NewAPI(userService *services.UsersService) *API {
	app := fiber.New()

	//Routes
	routes.APIRoutes(app)
	routes.UserRoutes(app, userService)

	PORT_API := os.Getenv("PORT_API")

	if PORT_API == "" {
		log.Fatal("the PORT_API is necesary, put it in the .env")
	}

	return &API{
		app:  app,
		PORT: PORT_API,
	}

}

func (api *API) Listen() {
	log.Fatal(api.app.Listen(":" + api.PORT))
}
