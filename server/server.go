package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.com/fibocloud/medtech/gin/utils"
)

// CORS cors enable
func CORS(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

// Start server
func Start() {
	// gin.DisableConsoleColor()
	//
	// Logging to a file.
	// f, _ := os.Create("log/gin.log")
	// gin.DefaultWriter = io.MultiWriter(f)
	//
	app := gin.Default()
	app.Use(CORS)
	r := Routers(app)
	utils.DocInit(app)
	r.Run(viper.GetString("server.PORT"))
}
