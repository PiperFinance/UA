package controllers

import (
	"fmt"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/schemas"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SignInUserNoSign(c *fiber.Ctx) error {
	var payload *schemas.SignInInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := schemas.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	users, add, user := []models.User{}, models.Address{}, models.User{}

	address, err := payload.Address.ETHAddress()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	if res := conf.DB.Model(&models.Address{}).Preload("Users").First(&add, "hash = ?", address.String()); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": res.Error.Error()})
	}

	//// Try session !!!
	//// .WithContext(c)
	//// if err := conf.DB.Debug().Model(&add).Association("Users").Find(&users); err != nil {
	//if err := conf.DB.Table("users").Preload("Users").Find(&address); err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "error": fmt.Sprintf("User Query: %s", err.Error.Error())})
	//}

	_ = users

	userFound := false
	for _, _user := range add.Users {
		err = bcrypt.CompareHashAndPassword([]byte(_user.Password), []byte(payload.Password))
		if err != nil {
			continue
		} else {
			user = *_user
			userFound = true
			break
		}
	}
	if !userFound {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
	}

	session, accessToken, err := GenAccessToken(user, nil)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating Access JWT Token failed: %v", err)})
	}
	refreshToken, err := GenRefreshToken(user, session)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating Refresh JWT Token failed: %v", err)})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   int(conf.Config.JwtMaxAge),
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(conf.Config.JwtMaxAge),
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"refreshToken": refreshToken,
		"accessToken":  accessToken,
	})
}
