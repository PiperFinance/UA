package main

import (
	"github.com/charmbracelet/log"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
		&models.SwapRequest{},
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
	app.Get("/address/", views.AllAddresses)
	app.Get("/users/", views.AllUsers)
	app.Get("/users/online", views.OnlineUsers)
	app.Get("/users/offline", views.OfflineUsers)
	app.Get("/", views.Accessible)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(conf.Config.JwtAccessSecret)},
	}))

	// Api with Needs Auth
	app.Post("/token/refresh", views.RefreshToken)
	app.Get("/token/validate", views.Validate)
	app.Get("/user/whoami", views.WhoAmI)
	app.Get("/user/address", views.GetUserAddresses)
	app.Post("/user/address", views.AddNewAddress)
	app.Delete("/user/address", views.RemoveAddress)

	app.Get("/user/tx", views.SwapReqHistory)
	app.Get("/user/tx/:uuid", views.SwapReqDet)
	app.Post("/user/tx", views.NewSwapReq)

	if err := app.Listen(conf.Config.ApiUrl); err != nil {
		log.Fatal(err)
	}
}
