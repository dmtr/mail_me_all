package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	appPrefix   string = "MAILME_APP"
	appHost     string = "127.0.0.1"
	appPort     int    = 8080
	fbProxyPort int    = 5000
	fbProxyHost string = "fbproxy"
	maxAge      int    = 60 * 60
)

// Config - app config
type Config struct {
	Host          string
	Port          int
	Debug         int
	Loglevel      log.Level
	DSN           string
	FbAppID       string
	FbRedirectURI string
	AppSecret     string
	FBProxyPort   int
	FBProxyHost   string
	Path          string
	Domain        string
	MaxAge        int
	Secure        bool
	HttpOnly      bool
	AuthKey       string
	EncryptKey    string
	Testing       bool
}

// GetConfig returns app config
func GetConfig() Config {
	viper.SetEnvPrefix(appPrefix)
	viper.SetDefault("HOST", appHost)
	viper.SetDefault("PORT", appPort)
	viper.SetDefault("DEBUG", 0)
	viper.SetDefault("Loglevel", "debug")
	viper.SetDefault("DSN", "")
	viper.SetDefault("FB_APP_ID", "")
	viper.SetDefault("FB_REDIRECT_URI", "")
	viper.SetDefault("FB_PROXY_PORT", fbProxyPort)
	viper.SetDefault("FB_PROXY_HOST", fbProxyHost)
	viper.SetDefault("PATH", "/")
	viper.SetDefault("DOMAIN", "localhost")
	viper.SetDefault("MAX_AGE", maxAge)
	viper.SetDefault("SECURE", 0)
	viper.SetDefault("HTTP_ONLY", 1)
	viper.SetDefault("TESTING", 0)
	viper.AutomaticEnv()

	loglevel, err := log.ParseLevel(viper.GetString("LOGLEVEL"))
	if err != nil {
		loglevel = log.ErrorLevel
	}

	conf := Config{
		Host:          viper.GetString("HOST"),
		Port:          viper.GetInt("PORT"),
		Debug:         viper.GetInt("DEBUG"),
		Loglevel:      loglevel,
		DSN:           viper.GetString("DSN"),
		FbAppID:       viper.GetString("FB_APP_ID"),
		FbRedirectURI: viper.GetString("FB_REDIRECT_URI"),
		AppSecret:     viper.GetString("app-secret"),
		FBProxyPort:   viper.GetInt("FB_PROXY_PORT"),
		FBProxyHost:   viper.GetString("FB_PROXY_HOST"),
		Path:          viper.GetString("PATH"),
		Domain:        viper.GetString("DOMAIN"),
		MaxAge:        viper.GetInt("MAX_AGE"),
		Secure:        viper.GetBool("SECURE"),
		HttpOnly:      viper.GetBool("HTTP_ONLY"),
		AuthKey:       viper.GetString("auth-key"),
		EncryptKey:    viper.GetString("encrypt-key"),
		Testing:       viper.GetBool("TESTING"),
	}

	return conf
}
