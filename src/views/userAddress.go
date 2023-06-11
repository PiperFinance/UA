package views

import (
	"strings"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
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

func AddNewAddress(c *fiber.Ctx) error {
	localUser := c.Locals("user").(*jwt.Token)
	claims := localUser.Claims.(jwt.MapClaims)
	var payload *schemas.MultiAddress

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	user := models.User{}
	if tx := conf.DB.First(&user, "uuid = ?", claims["sub"]); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": tx.Error.Error()})
	}

	for i := range payload.Addresses {
		if tx := conf.DB.Create(&payload.Addresses[i]); tx.Error != nil {
			if strings.Contains(tx.Error.Error(), "duplicate key") {
				conf.DB.Find(&payload.Addresses[i], "hash = ?", payload.Addresses[i].Hash)
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": tx.Error.Error()})
			}
		}
	}
	user.Addresses = payload.Addresses
	if tx := conf.DB.Save(&user); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": tx.Error.Error()})
	}
	return c.SendStatus(201)
}

func RemoveAddress(c *fiber.Ctx) error {
	// TODO .... This Does not work
	localUser := c.Locals("user").(*jwt.Token)
	claims := localUser.Claims.(jwt.MapClaims)
	var payload *schemas.MultiAddress

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	user := models.User{}
	if tx := conf.DB.Preload("Addresses").First(&user, "uuid = ?", claims["sub"]); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": tx.Error.Error()})
	}
	r := make([]*models.Address, 0)
	for _, selectedAdd := range payload.Addresses {
		for _, add := range user.Addresses {
			if add.Hash == selectedAdd.Hash {
				continue
			}
			r = append(r, add)
		}
	}
	user.Addresses = r
	if tx := conf.DB.Save(&user); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": tx.Error.Error()})
	}
	return c.SendStatus(204)
}

func GetUserAddresses(c *fiber.Ctx) error {
	localUser := c.Locals("user").(*jwt.Token)
	claims := localUser.Claims.(jwt.MapClaims)
	user := models.User{}

	if tx := conf.DB.Preload("Addresses").First(&user, "uuid = ?", claims["sub"]); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": tx.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "OK", "addresses": user.Addresses,
	})
}
