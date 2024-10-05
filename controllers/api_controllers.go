package controllers

import "github.com/gofiber/fiber/v2"

// PingAPI controll the endpoint for make an ping
// it is usefull for know if the api its working
func PingAPI(c *fiber.Ctx) error {
	return c.SendStatus(204)
}
