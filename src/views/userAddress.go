package views

import (
	"github.com/PiperFinance/UA/src/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func UpdateUserAddress(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(schemas.TokenClaim)
	claims.GetAddresses()
	return nil
}

func GetUserAddresses(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(schemas.TokenClaim)
	adds, err := claims.GetAddresses()
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status": "fail", "error": err,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success", "adddresses": adds,
		})
	}
}
