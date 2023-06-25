package views

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/schemas"
)

func NewSwapReq(c *fiber.Ctx) error {
	var payload *models.SwapRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(422).JSON(fiber.Map{"err": err})
	}
	errors := schemas.ValidateStruct(payload)
	if errors != nil {
		return c.Status(422).JSON(fiber.Map{"err": errors})
	}

	localUser := c.Locals("user").(*jwt.Token)
	claims := localUser.Claims.(jwt.MapClaims)
	user := models.User{}
	if tx := conf.DB.First(&user, "uuid = ?", claims["sub"]); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": tx.Error.Error()})
	} else {
		payload.User = user
		payload.UserUUID = user.UUID
	}

	if tx := conf.DB.Debug().Save(&payload); tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"err": tx.Error.Error()})
	} else {
		return c.Status(200).JSON(fiber.Map{"result": payload})
	}
}

func SwapReqHistory(c *fiber.Ctx) error {
	localUser := c.Locals("user").(*jwt.Token)
	claims := localUser.Claims.(jwt.MapClaims)
	user := models.User{}
	if tx := conf.DB.First(&user, "uuid = ?", claims["sub"]); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": tx.Error.Error()})
	} else {
		res := make([]models.SwapRequest, 0)
		if tx := conf.DB.Preload("User").Preload("Address").Find(&res, "user_uuid = ?", user.UUID); tx.Error != nil {
			return c.Status(422).JSON(fiber.Map{"err": tx.Error.Error()})
		}
		return c.Status(200).JSON(fiber.Map{"result": res})
	}
}

func SwapReqDet(c *fiber.Ctx) error {
	// TODO - check if user has access
	uid := c.Params("uuid")
	rq := models.SwapRequest{}
	if tx := conf.DB.Preload("User").Preload("Address").First(&rq, "uuid = ?", uid); tx.Error != nil {
		return c.Status(422).JSON(fiber.Map{"err": tx.Error.Error()})
	} else {
		return c.Status(200).JSON(fiber.Map{"result": rq})
	}
}
