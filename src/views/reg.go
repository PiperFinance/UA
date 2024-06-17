package views

import (
	"github.com/gofiber/fiber/v2"

	"github.com/PiperFinance/UA/src/controllers"
)

func SignUpUser(c *fiber.Ctx) error {
	return controllers.SignUpUser(c)
}

func SignUpAndSignInUser(c *fiber.Ctx) error {
	if err := controllers.SignUpUser(c); err != nil {
		return err
	}
	return controllers.SignInUser(c)
}

func SignUpAndSignInUserNoSign(c *fiber.Ctx) error {
	controllers.SignUpUserNoSign(c)
	// if r != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": r.Error()})
	// }
	return controllers.SignInUserNoSign(c)
}
