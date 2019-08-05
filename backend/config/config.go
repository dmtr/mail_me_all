package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	appPrefix string = "MAILME_APP"
	appHost   string = "127.0.0.1"
	appPort   int    = 8080
)

// Config - app config
type Config struct {
	Host string
	Port int
}

// GetConfig returns app config
func GetConfig() Config {
	viper.SetEnvPrefix(appPrefix)
	viper.SetDefault("HOST", appHost)
	viper.SetDefault("PORT", appPort)
	viper.AutomaticEnv()

	conf := Config{
		Host: viper.GetString("HOST"),
		Port: viper.GetInt("PORT"),
	}

	log.Println("Config loaded")
	return conf
}
