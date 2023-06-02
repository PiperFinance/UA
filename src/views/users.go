package views

import (
	"time"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/gofiber/fiber/v2"
)

func OnlineUsers(c *fiber.Ctx) error {
	var sessions []*models.Session
	if res := conf.DB.Model(&models.Session{}).Preload("User").Find(&sessions, "expires_at >= ?", time.Now().Unix()); res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "error": res.Error.Error(),
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "OK", "msg": sessions,
		})
	}
}

func AllUsers(c *fiber.Ctx) error {
	var sessions []*models.Session
	if res := conf.DB.Model(&models.Session{}).Preload("User").Preload("Addresses").Preload("Device").Find(&sessions); res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "error": res.Error.Error(),
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "OK", "msg": sessions,
		})
	}
}

func OfflineUsers(c *fiber.Ctx) error {
	var sessions []*models.Session
	if res := conf.DB.Model(&models.Session{}).Preload("User").Preload("Addresses").Preload("Device").Find(&sessions, "expires_at < ?", time.Now().Unix()); res.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "error": res.Error.Error(),
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "OK", "msg": sessions,
		})
	}
}
