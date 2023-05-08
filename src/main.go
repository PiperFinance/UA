package main

import (
	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/views"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func init() {
	if err := conf.LoadConfig("."); err != nil {
		log.Fatal(err)
	}
	if err := conf.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	if err := conf.DB.AutoMigrate(
		&models.User{},
		&models.Address{},
		&models.Device{},
	); err != nil {
		log.Fatal(err)
	}
}

func main() {

	app := fiber.New()
	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:        []byte(conf.Config.JwtAccessSecret),
		KeyRefreshTimeout: &conf.Config.JwtRefreshExpiresIn}))
	//Api
	app.Get("/api/healthchecker", views.HealthCheck)
	//User
	app.Post("/login", views.Login)
	app.Post("/signup", views.SignUpUser)
	app.Post("/refresh", views.RefreshToken)

	// Unauthenticated route
	app.Get("/", views.Accessible)
	// Restricted Routes
	app.Get("/restricted", views.Restricted)
	app.Post("/", views.BodyParserSample)

	if err := app.Listen(":4500"); err != nil {
		log.Fatal(err)
	}

}
