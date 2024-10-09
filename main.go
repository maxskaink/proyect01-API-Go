package main

import (
	api "github.com/maxskaink/proyect01-api-go/API"
	"github.com/maxskaink/proyect01-api-go/dataccess/repositories"
	"github.com/maxskaink/proyect01-api-go/services"
	"github.com/maxskaink/proyect01-api-go/utils"
)

// main entry point for mi API rest aplication
// get de eviroment variables, connect with de database
// get de routes for de api and init the litening of the API
func main() {
	utils.LoadENV()

	userRepository := repositories.NewUserMongoRepository()
	userService := services.NewUsersService(userRepository)

	appAPI := api.NewAPI(userService)
	defer userRepository.CloseClient()

	appAPI.Listen()

}
