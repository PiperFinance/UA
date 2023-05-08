package api

import (
	"bytes"
	"fmt"
	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type AuthReq struct {
	ChainId   int32     `json:"chainId" validate:"required"`
	Hex       string    `validate:"required" `
	UserUUID  uuid.UUID `json:"UID" validate:"uuid"`
	SignedMsg []byte
}

func GetOrCreateUser(ctx fiber.Ctx, authReq AuthReq) (*models.User, error) {
	add := models.Address{}
	if err := conf.DB.FirstOrCreate(&add, "hash = ?", authReq.Hex).Error; err != nil {
		return nil, err
	}
	user := models.User{}
	if bytes.Equal(authReq.UserUUID[:], (&uuid.UUID{})[:]) {
		conf.DB.FirstOrCreate(&user, "uuid = ?", authReq.UserUUID)
	} else {
		conf.DB.FirstOrCreate(&user, "uuid = ?", authReq.UserUUID)
	}
	user.Addresses = append(user.Addresses, add)
	return &user, nil
}

func AuthUser(c *fiber.Ctx) error {
	//c.
	auth := new(AuthReq)
	if err := c.BodyParser(&auth); err != nil {
		return err
	}
	ok, err := auth.Address.VerifySignedBytes(auth.Address.String(), auth.SignedMsg)
	if err != nil {
		return nil
	}
	if ok {
		user, err := GetOrCreateUser(*c, *auth)
		if err != nil {
			return err
		}
		fmt.Println(user)
	}
	return err
}

func OnlineUser(c *fiber.Ctx) error {
	startingTime := time.Now()
	user := models.User{}
	//c.ShouldBindBodyWith(&user)
	conf.DB.Create(&user)
	_ = startingTime
}

//func RefreshhUser()

func AllUsers(c *gin.Context) {
	var users []models.User
	if err := conf.DB.Find(&users).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, users)
	}
}

func LastUser(c *gin.Context) {
	user := models.User{}
	conf.DB.First(&user)
	c.JSON(http.StatusOK, user)
}
