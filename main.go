package main

import (
	"flag"
	"fmt"
	"os"

	config "gitlab.com/fibocloud/medtech/gin/config"
	server "gitlab.com/fibocloud/medtech/gin/server"
)

// @title MedTech
// @version 1.0
// @description MedTech

// @contact.name FIBO CLOUD
// @contact.url https://www.facebook.com/fibocloud/
// @contact.email we@fibo.cloud

// @license.name MIT License
// @license.url https://en.wikipedia.org/wiki/MIT_License

// @host 127.0.0.1:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey Token
// @in header
// @name Authorization
func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	server.Start()
}
