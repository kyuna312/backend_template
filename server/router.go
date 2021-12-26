package server

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/fibocloud/medtech/gin/controllers"
)

// Routers ...
func Routers(app *gin.Engine) *gin.Engine {
	api := app.Group("/api/v1")
	controllers.Init(api)
	return app
}
