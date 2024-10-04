package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/maxskaink/proyect01-api-go/models"
	"github.com/maxskaink/proyect01-api-go/models/dto"
	"github.com/maxskaink/proyect01-api-go/services"
)

func CreateUser(c *fiber.Ctx) error {
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

	userCreated, err := services.CreateUser(newUser)
	if err != nil {
		return c.Status(500).JSON(dto.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}
	return c.Status(201).JSON(userCreated)
}

func GetAllUsers(c *fiber.Ctx) error {
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
	totalUsers, err := services.GetTotalUsers()

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

	users, err := services.GetAllUsers(page, per_page)
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

func LogInUser(c *fiber.Ctx) error {
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

	token, err := services.LogInUser(credential.Email, credential.Password)

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

func GetUserbyID(c *fiber.Ctx) error {
	idToSearch := c.Params("id")
	if idToSearch == "" {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	foundUser, err := services.GetUserByID(idToSearch)

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

func UpdateUser(c *fiber.Ctx) error {
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

	userDB, err := services.GetUserByID(idStr)

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

	oldUser, err := services.ReplaceUser(newUser, idStr)

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
