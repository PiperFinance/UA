package main

import (
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
	_ "github.com/joho/godotenv/autoload"

	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/views"
)

func init() {
	conf.LoadConfig()
	conf.ConnectMongo()
	conf.LoadLogger()
	conf.ConnectDB()
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
	app.Post("/SignUpSignInNoSign", views.SignUpAndSignInUserNoSign)
	app.Post("/refresh", views.RefreshToken)
	app.Get("/address/", views.AllAddresses)
	app.Get("/users/", views.AllUsers)
	app.Get("/users/online", views.OnlineUsers)
	app.Get("/users/offline", views.OfflineUsers)
	app.Get("/", views.Accessible)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:        []byte(conf.Config.JwtAccessSecret),
		KeyRefreshTimeout: &conf.Config.JwtRefreshExpiresIn,
	}))
	// Api with Needs Auth

	app.Get("/validate", views.Validate)
	app.Get("/whoami", views.WhoAmI)
	app.Get("/user/address", views.GetUserAddresses)
	app.Post("/user/address", views.AddNewAddress)
	app.Delete("/user/address", views.RemoveAddress)

	if err := app.Listen(conf.Config.ApiUrl); err != nil {
		log.Fatal(err)
	}
}
