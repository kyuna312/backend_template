package utils

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger documentation register
	_ "gitlab.com/fibocloud/medtech/gin/docs"
)

// DocInit swagger documentation init
func DocInit(app *gin.Engine) {
	url := ginSwagger.URL("http://103.50.204.55:8080/swagger/doc.json")
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
