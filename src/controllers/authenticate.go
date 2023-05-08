package controllers

import (
	"fmt"
	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

func GenRefreshToken(user models.User) (string, error) {
	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = user.UUID
	claims["exp"] = now.Add(conf.Config.JwtRefreshExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	return tokenByte.SignedString([]byte(conf.Config.JwtRefreshSecret))

}
func GenAccessToken(user models.User) (string, error) {

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["adds"] = user.Add2Str()
	claims["sub"] = user.UUID
	claims["exp"] = now.Add(conf.Config.JwtAccessExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	return tokenByte.SignedString([]byte(conf.Config.JwtAccessSecret))

}

func RefreshToken(refreshToken *jwt.Token) (string, string, error) {
	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", jwt.ErrInvalidKey
	}
	userUUID := claims["sub"].(string)
	user := models.User{}
	if res := conf.DB.First(&user, "uuid = ?", userUUID); res.Error != nil {
		return "", "", res.Error
	}
	refT, refErr := GenRefreshToken(user)
	if refErr != nil {
		return "", "", refErr
	}
	accT, accErr := GenAccessToken(user)
	if refErr != nil {
		return "", "", accErr
	}
	user.LastAccess = time.Now().UTC()
	return accT, refT, nil
}

func SignInUser(c *fiber.Ctx) error {
	var payload *schemas.SignInInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := schemas.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	var user models.User
	// TODO - GET User By Address M2M
	address, err := payload.Address.ETHAddress()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	result := conf.DB.First(&user, "address = ?", strings.ToLower(address.String()))

	if result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
	}

	refreshToken, err := GenRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating Refresh JWT Token failed: %v", err)})
	}
	accessToken, err := GenAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating Access JWT Token failed: %v", err)})
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
		"accessToken":  accessToken})
}
