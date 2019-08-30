package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	appPrefix string = "MAILME_APP"
	appHost   string = "127.0.0.1"
	appPort   int    = 8080
)

// Config - app config
type Config struct {
	Host     string
	Port     int
	Debug    int
	Loglevel log.Level
	DSN      string
}

// GetConfig returns app config
func GetConfig() Config {
	viper.SetEnvPrefix(appPrefix)
	viper.SetDefault("HOST", appHost)
	viper.SetDefault("PORT", appPort)
	viper.SetDefault("DEBUG", 0)
	viper.SetDefault("Loglevel", "debug")
	viper.SetDefault("DSN", "")
	viper.AutomaticEnv()

	loglevel, err := log.ParseLevel(viper.GetString("LOGLEVEL"))
	if err != nil {
		loglevel = log.ErrorLevel
	}

	conf := Config{
		Host:     viper.GetString("HOST"),
		Port:     viper.GetInt("PORT"),
		Debug:    viper.GetInt("DEBUG"),
		Loglevel: loglevel,
		DSN:      viper.GetString("DSN"),
	}

	return conf
}
