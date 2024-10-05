package auth

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// UserAuth Its a middleware for validate if the actual http query
// have the jwt to make an secure secion
func UserAuth() fiber.Handler {
	key := os.Getenv("JWT_SECRET")
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(key)},
		SuccessHandler: func(c *fiber.Ctx) error {
			user, ok := c.Locals("user").(*jwt.Token)
			if !ok {
				return c.Status(500).JSON(fiber.Map{
					"message": "Error getting user token",
				})
			}

			// Obtener los claims del token
			claims, ok := user.Claims.(jwt.MapClaims)
			if !ok {
				return c.Status(500).JSON(fiber.Map{
					"message": "Error getting user claims",
				})
			}

			// Crear una nueva instancia de Credential y asignar los valores de los claims

			email, emailOk := claims["email"].(string)
			name, nameOk := claims["name"].(string)

			if !emailOk || !nameOk {
				return c.Status(500).JSON(fiber.Map{
					"message": "Error getting email or password from claims",
				})
			}

			c.Locals("email", email)
			c.Locals("name", name)
			return c.Next()
		},
	})
}
