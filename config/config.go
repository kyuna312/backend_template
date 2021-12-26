package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Init config file
func Init(env string) {
	var err error
	fmt.Println(env)
	viper.SetConfigType("yaml")
	viper.SetConfigName(env)
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("config/")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
}
