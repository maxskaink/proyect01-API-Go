package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxskaink/proyect01-api-go/models"
	"github.com/maxskaink/proyect01-api-go/models/dto"
	"github.com/maxskaink/proyect01-api-go/services"
)

type ClientControllers struct {
	UserService *services.UsersService
}

func NewClientControllers(service *services.UsersService) *ClientControllers {
	return &ClientControllers{
		UserService: service,
	}
}

// CreateUser handle the endpoint for create and user
// validate the information and use the service to create
// and save in the database de information.
func (clientC *ClientControllers) CreateUser(c *fiber.Ctx) error {
	newUser := new(models.User)
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(400).JSON(dto.Error{
			Message:   "Error parsing body",
			Status:    400,
			TypeError: "Invalid sintaxis",
		})
	}

	if err := newUser.ValidateToCreate(); err != nil {
		return c.Status(400).JSON(dto.Error{
			Message:   err.Error(),
			Status:    400,
			TypeError: "Invalid data",
		})
	}

	userCreated, err := clientC.UserService.CreateUser(newUser)
	if err != nil {
		return c.Status(500).JSON(dto.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}
	return c.Status(201).JSON(userCreated)
}

// GetAllUsers handle the endpoint for get some users
// with the specifications of the user
func (clientC *ClientControllers) GetAllUsers(c *fiber.Ctx) error {
	//Sacmos la info de querys y validamos
	querys := c.Queries()

	page, err := strconv.Atoi(querys["page"])
	if err != nil {
		page = 1
	}

	per_page, err := strconv.Atoi(querys["per_page"])
	if err != nil {
		per_page = 10
	}
	if per_page <= 0 || page <= 0 {
		return c.Status(400).JSON(dto.Error{
			Message:   "Invalid query params, per_page and page must be greater than 0",
			Status:    400,
			TypeError: "Invalid data",
		})
	}

	//Sacmos los usuarios de la base de datos y validamos info
	totalUsers, err := clientC.UserService.GetTotalUsers()

	if err != nil {
		return c.Status(500).JSON(dto.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	} else if totalUsers == 0 {
		return c.Status(404).JSON(dto.Error{
			Message:   "No users found",
			Status:    404,
			TypeError: "Not found",
		})
	}

	totalPages := int(totalUsers/per_page) + 1

	if page > totalPages {
		return c.Status(404).JSON(dto.Error{
			Message:   "Page not found",
			Status:    404,
			TypeError: "Not found",
		})
	}

	users, err := clientC.UserService.GetAllUsers(page, per_page)
	if err != nil {
		return c.Status(500).JSON(dto.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}

	//Cremos la respuesta
	response := dto.ResponseAllUsers{
		Data:           users,
		Page:           page,
		TotalPages:     totalPages,
		TotalUsersPage: len(users),
	}

	return c.Status(200).JSON(response)
}

// LogInUser handle the endpoint for login an user
// returning in a json a JWT
func (clientC *ClientControllers) LogInUser(c *fiber.Ctx) error {
	credential := new(dto.Credential)

	if err := c.BodyParser(credential); err != nil {
		return c.Status(400).JSON(dto.Error{
			Message:   "Error parsing body",
			Status:    400,
			TypeError: "Invalid sintaxis",
		})
	}
	if credential.Email == "" || credential.Password == "" {
		return c.Status(400).JSON(dto.Error{
			Message:   "Email and password are required",
			Status:    400,
			TypeError: "Invalid data",
		})
	}

	token, err := clientC.UserService.LogInUser(credential.Email, credential.Password)

	if err != nil {
		return c.Status(401).JSON(dto.Error{
			Message:   err.Error(),
			Status:    401,
			TypeError: "Unauthorized",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"token": token,
	})
}

// GetUserbyID hanlde the endpoint for get all the information
// of an specific user, it mus be have and jwt
func (clientC *ClientControllers) GetUserbyID(c *fiber.Ctx) error {
	idToSearch := c.Params("id")
	if idToSearch == "" {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	foundUser, err := clientC.UserService.GetUserByID(idToSearch)

	if err != nil {
		return c.Status(500).JSON(dto.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}

	if foundUser.Email != c.Locals("email") {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Error{
			Message:   "Not auhorized, to get information of this user",
			Status:    fiber.StatusUnauthorized,
			TypeError: "Unahorized",
		})
	}

	return c.Status(200).JSON(foundUser)
}

// UpdateUser handle the nedpoint for update all the information
// of the user, the user must be have an JWT
func (clientC *ClientControllers) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	newUser := new(models.User)

	err := c.BodyParser(newUser)

	if err != nil {
		return c.Status(400).JSON(dto.Error{
			Message:   err.Error(),
			Status:    400,
			TypeError: "Invalid Body",
		})
	}

	userDB, err := clientC.UserService.GetUserByID(idStr)

	if err != nil {
		return c.Status(500).JSON(dto.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}

	if userDB.Email != c.Locals("email") {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Error{
			Message:   "Not authorized for the update",
			Status:    fiber.StatusUnauthorized,
			TypeError: "Client error",
		})
	}

	oldUser, err := clientC.UserService.ReplaceUser(newUser, idStr)

	if err != nil {
		return c.Status(500).JSON(dto.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "ok",
		"oldUser": oldUser,
	})

}

// PatchUser handle the endpoint for update some information of the user
// the user must have a JWT, and the information to update must be in the body
func (clientC *ClientControllers) PatchUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	newUser := new(models.User)

	if err := c.BodyParser(newUser); err != nil {
		return c.Status(400).JSON(dto.Error{
			Message:   err.Error(),
			Status:    400,
			TypeError: "Invalid Body",
		})
	}

	userDB, err := clientC.UserService.GetUserByID(idStr)

	if err != nil {
		return c.Status(500).JSON(dto.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}

	if userDB.Email != c.Locals("email") {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.Error{
			Message:   "Not authorized for the update",
			Status:    fiber.StatusUnauthorized,
			TypeError: "Client error",
		})
	}

	oldUser, err := clientC.UserService.UpdateUserById(newUser, idStr)

	if err != nil {
		return c.Status(500).JSON(dto.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}
	oldUser.Password = ""

	return c.Status(200).JSON(oldUser)

}

// DeleteUser handle the enpoint for update the state of isActive to false
// if its any error it will return error
func (clientC *ClientControllers) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")

	userDB, err := clientC.UserService.GetUserByID(idStr)
	if err != nil {
		return c.Status(500).JSON(dto.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error getting user",
		})
	}

	if userDB.Email != c.Locals("email") {
		return c.Status(401).JSON(dto.Error{
			Message:   "You are not authorized",
			Status:    401,
			TypeError: "Unathorized for this action",
		})
	}

	deletedUser, err := clientC.UserService.DeleteUserById(idStr)

	if err != nil {
		return c.Status(500).JSON(dto.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}

	return c.Status(200).JSON(deletedUser)
}
