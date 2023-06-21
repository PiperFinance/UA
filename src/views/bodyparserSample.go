package views

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
)

func WhoAmI(c *fiber.Ctx) error {
	localUser := c.Locals("user").(*jwt.Token)
	claims := localUser.Claims.(jwt.MapClaims)
	user := models.User{}
	if tx := conf.DB.First(&user, "uuid = ?", claims["sub"]); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": tx.Error.Error()})
	}

	return c.JSON(fiber.Map{"claims": claims, "user": user})
}
