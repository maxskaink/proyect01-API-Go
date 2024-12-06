package controllers

import (
	"github.com/gofiber/fiber/v2"
	custom_errors "github.com/maxskaink/proyect01-api-go/errors"
	"github.com/maxskaink/proyect01-api-go/models/dto"
)

func ResponseError(err error, c *fiber.Ctx) error {
	if instanceError, ok := err.(*custom_errors.InvalidFormat); ok {
		return c.Status(400).JSON(dto.Error{
			Message:   instanceError.Error(),
			Status:    400,
			TypeError: "Invalid data",
		})
	}

	if instanceError, ok := err.(*custom_errors.UnAuthorized); ok {
		return c.Status(401).JSON(dto.Error{
			Message:   instanceError.Error(),
			Status:    401,
			TypeError: "Unauthorized",
		})
	}

	if instanceError, ok := err.(*custom_errors.NotFound); ok {
		return c.Status(404).JSON(dto.Error{
			Message:   instanceError.Error(),
			Status:    404,
			TypeError: "Not Found",
		})
	}

	if instanceError, ok := err.(*custom_errors.DuplicateInformation); ok {
		return c.Status(409).JSON(dto.Error{
			Message:   instanceError.Error(),
			Status:    409,
			TypeError: "Duplicate Information",
		})
	}

	if instanceError, ok := err.(*custom_errors.InternalError); ok {
		return c.Status(500).JSON(dto.Error{
			Message:   instanceError.Error(),
			Status:    500,
			TypeError: "Internal error",
		})
	}

	return nil
}
