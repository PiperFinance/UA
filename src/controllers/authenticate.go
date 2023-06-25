package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/schemas"
)

func GenRefreshToken(user models.User) (string, error) {
	now := time.Now().UTC()
	expiry := now.Add(conf.Config.JwtRefreshExpiresIn).Unix()
	claims := jwt.MapClaims{
		"sub": *user.UUID,
		"exp": expiry,
		"iat": now.Unix(),
		"nbf": now.Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(conf.Config.JwtRefreshSecret))
}

func GenAccessToken(user models.User, PrevSession *models.Session) (string, error) {
	now := time.Now().UTC()
	expiry := now.Add(conf.Config.JwtAccessExpiresIn).Unix()
	claims := jwt.MapClaims{
		"sub": *user.UUID,
		"exp": expiry,
		"iat": now.Unix(),
		"nbf": now.Unix(),
	}
	session := models.Session{
		ExpiresAt: expiry,
		UserRefer: user.UUID,
		User:      user,
	}
	if PrevSession != nil {
		session.UUID = PrevSession.UUID
		if res := conf.DB.Save(&session); res.Error != nil {
			return "", res.Error
		}
	} else {
		if res := conf.DB.FirstOrCreate(&session); res.Error != nil {
			return "", res.Error
		}
		session.ExpiresAt = expiry
		if res := conf.DB.Save(&session); res.Error != nil {
			return "", res.Error
		}

		claims["suid"] = *session.UUID
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(conf.Config.JwtAccessSecret))
}

func RefreshToken(refreshToken *jwt.Token) (string, string, error) {
	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", jwt.ErrInvalidKey
	}
	userUUID := claims["sub"].(string)
	sessionUUID := claims["suid"].(string)

	user, session := models.User{}, models.Session{}
	if res := conf.DB.First(&user, "uuid = ?", userUUID); res.Error != nil {
		return "", "", res.Error
	}
	if res := conf.DB.First(&session, "uuid = ?", sessionUUID); res.Error != nil {
		return "", "", res.Error
	}
	accT, accErr := GenAccessToken(user, &session)
	if accErr != nil {
		return "", "", accErr
	}
	refT, refErr := GenRefreshToken(user)
	if refErr != nil {
		return "", "", refErr
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

	accessToken, err := GenAccessToken(user, nil)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating Access JWT Token failed: %v", err)})
	}
	refreshToken, err := GenRefreshToken(user)
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
