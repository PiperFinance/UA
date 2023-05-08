package views

import (
	"github.com/PiperFinance/UA/src/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func Restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}

func RefreshToken(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	accessToken, refreshToken, err := controllers.RefreshToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":       "success",
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}

func Login(c *fiber.Ctx) error {
	return controllers.SignInUser(c)
}
