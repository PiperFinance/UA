package controllers

import (
	"fmt"
	"strings"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/jobs"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/schemas"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SignUpUserNoSign(c *fiber.Ctx) error {
	var payload *schemas.SignUpInputNoSign

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	errors := schemas.ValidateStruct(payload)
	if errors != nil {
		return fmt.Errorf("%+v", errors)
	}

	var err error

	if res := conf.DB.FirstOrCreate(&payload.Address); res.Error != nil {
		return res.Error
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := models.User{
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: string(hashedPassword),
	}
	newUser.Addresses = []*models.Address{&payload.Address}

	result := conf.DB.Create(&newUser)
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return fmt.Errorf("user with that email already exists")
	} else if result.Error != nil {
		return fmt.Errorf("something bad happened , err : %+v", result.Error)
	}

	o := (jobs.SyncAddress{Address: newUser.Addresses[0]})
	go func(o jobs.SyncAddress) {
		if err := o.ExecuteAll(); err != nil {
			conf.Logger.Error(err)
			fmt.Println(err)

		}
	}(o)

	return nil
}
