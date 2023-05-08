package api

import "github.com/gin-gonic/gin"

var (
	app *gin.Engine
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	app = gin.Default()
	app.GET("/", hello)
	return app
}
