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
