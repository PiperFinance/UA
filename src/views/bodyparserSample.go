package views

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthReq struct {
	ChainId    int32     `json:"chainId" validate:"required"`
	EthAddress string    `json:"ethAdd" validate:"required,eth_addr" `
	UserUUID   uuid.UUID `json:"UID" validate:"uuid"`
	SignedMsg  []byte    `json:"signedMsg" validate:"hexadecimal"`
}

func BodyParserSample(c *fiber.Ctx) error {
	authReq := AuthReq{}
	if err := c.BodyParser(&authReq); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(authReq)
}
