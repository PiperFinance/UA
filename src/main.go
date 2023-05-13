package main

import (
	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/views"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
		&models.Session{},
	); err != nil {
		log.Fatal(err)
	}
}

func main() {

	app := fiber.New()
	// Initialize default config
	app.Use(cors.New())

	// No Auth
	app.Get("/api/healthchecker", views.HealthCheck)
	app.Post("/login", views.Login)
	app.Post("/signup", views.SignUpUser)
	app.Post("/SignUpSignIn", views.SignUpAndSignInUser)
	app.Post("/refresh", views.RefreshToken)
	app.Get("/users", views.OnlineUsers)
	app.Get("/", views.Accessible)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:        []byte(conf.Config.JwtAccessSecret),
		KeyRefreshTimeout: &conf.Config.JwtRefreshExpiresIn}))
	// Api with Needs Auth

	// Unauthenticated route
	// Restricted Routes
	app.Get("/restricted", views.Restricted)
	app.Post("/", views.BodyParserSample)

	if err := app.Listen(":4500"); err != nil {
		log.Fatal(err)
	}

}
