package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxskaink/proyect01-api-go/models"
	"github.com/maxskaink/proyect01-api-go/services"
)

func CreateUser(c *fiber.Ctx) error {
	newUser := new(models.User)
	if err := c.BodyParser(newUser); err != nil {
		return c.Status(400).JSON(models.Error{
			Message:   "Error parsing body",
			Status:    400,
			TypeError: "Invalid sintaxis",
		})
	}

	if err := newUser.ValidateToCreate(); err != nil {
		return c.Status(400).JSON(models.Error{
			Message:   err.Error(),
			Status:    400,
			TypeError: "Invalid data",
		})
	}

	userCreated, err := services.CreateUser(newUser)
	if err != nil {
		return c.Status(500).JSON(models.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}
	return c.Status(201).JSON(userCreated)
}

func GetAllUsers(c *fiber.Ctx) error {
	users, err := services.GetAllUsers()
	if err != nil {
		return c.Status(500).JSON(models.Error{
			Message:   err.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}
	return c.Status(200).JSON(users)
}
