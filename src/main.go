package main

import (
	"github.com/charmbracelet/log"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hibiken/asynq"
	_ "github.com/joho/godotenv/autoload"

	"github.com/PiperFinance/UA/src/bg"
	"github.com/PiperFinance/UA/src/bg/handlers"
	"github.com/PiperFinance/UA/src/conf"
	"github.com/PiperFinance/UA/src/models"
	"github.com/PiperFinance/UA/src/views"
)

func init() {
	conf.LoadConfig()
	conf.ConnectMongo()
	conf.LoadLogger()
	conf.LoadRedis()
	conf.ConnectDB()
	conf.LoadCronTab()
	conf.LoadQueue()
	go conf.RunWorker(xHandlers())
	go conf.RunScheduler(xSchedules())
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

func xHandlers() []conf.MuxHandler {
	return []conf.MuxHandler{
		{Key: bg.SyncNTScheduleTaskKey, Handler: handlers.SyncNTFsScheduleTaskHandler, Q: asynq.Queue(conf.UASyncNTQ)},
		{Key: bg.SyncNTTaskKey, Handler: handlers.SyncNTFsTaskHandler, Q: asynq.Queue(conf.UASyncNTQ)},
		{Key: bg.SyncTHScheduleTaskKey, Handler: handlers.SyncTrxScheduleTaskHandler, Q: asynq.Queue(conf.UASyncTHQ)},
		{Key: bg.SyncTHTaskKey, Handler: handlers.SyncTrxTaskHandler, Q: asynq.Queue(conf.UASyncTHQ)},
		{Key: bg.SyncPairBalTaskKey, Handler: handlers.PairBalTaskHandler, Q: asynq.Queue(conf.UASyncBalQ)},
		{Key: bg.SyncTokenBalTaskKey, Handler: handlers.TokenBalTaskHandler, Q: asynq.Queue(conf.UASyncBalQ)},
	}
}

func xSchedules() []conf.QueueSchedules {
	return []conf.QueueSchedules{
		{Cron: "@every 3m", Q: asynq.Queue(conf.UASyncTHQ), Timeout: conf.Config.THSaveTimeout, Key: bg.SyncTHScheduleTaskKey},
		{Cron: "@every 1m", Q: asynq.Queue(conf.UASyncNTQ), Timeout: conf.Config.NTSaveTimeout, Key: bg.SyncNTScheduleTaskKey},
	}
}

func main() {
	app := fiber.New()
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
