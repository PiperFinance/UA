package views

import (
	"github.com/PiperFinance/UA/src/controllers"
	"github.com/gofiber/fiber/v2"
)

func SignUpUser(c *fiber.Ctx) error {
	return controllers.SignUpUser(c)
}

func SignUpAndSignInUser(c *fiber.Ctx) error {

	controllers.SignUpUser(c)
	return controllers.SignInUser(c)
}
