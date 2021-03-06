package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	appPrefix   string = "MAILME_APP"
	appHost     string = "127.0.0.1"
	appPort     int    = 8080
	twProxyPort int    = 5000
	twProxyHost string = "twproxy"
	maxAge      int    = 60 * 60
)

// Config - app config
type Config struct {
	Host             string
	Port             int
	Debug            int
	Loglevel         log.Level
	DSN              string
	TwProxyPort      int
	TwProxyHost      string
	Path             string
	Domain           string
	MaxAge           int
	Secure           bool
	HttpOnly         bool
	AuthKey          string
	EncryptKey       string
	Testing          bool
	TwConsumerKey    string
	TwConsumerSecret string
	TwCallbackURL    string
	TemplatePath     string
	MgDomain         string
	MgAPIKEY         string
	From             string
	PemFile          string
	KeyFile          string
	TweetTTL         int
}

// GetConfig returns app config
func GetConfig() Config {
	viper.SetEnvPrefix(appPrefix)
	viper.SetDefault("HOST", appHost)
	viper.SetDefault("PORT", appPort)
	viper.SetDefault("DEBUG", 0)
	viper.SetDefault("Loglevel", "debug")
	viper.SetDefault("DSN", "")
	viper.SetDefault("TW_PROXY_PORT", twProxyPort)
	viper.SetDefault("TW_PROXY_HOST", twProxyHost)
	viper.SetDefault("PATH", "/")
	viper.SetDefault("DOMAIN", "localhost")
	viper.SetDefault("MAX_AGE", maxAge)
	viper.SetDefault("SECURE", 0)
	viper.SetDefault("HTTP_ONLY", 1)
	viper.SetDefault("TESTING", 0)
	viper.SetDefault("TW_CALLBACK_URL", "")
	viper.SetDefault("TEMPLATE_PATH", "/app/templates/")
	viper.SetDefault("MG_DOMAIN", "")
	viper.SetDefault("MG_APIKEY", "")
	viper.SetDefault("FROM", "")
	viper.SetDefault("PEM_FILE", "")
	viper.SetDefault("KEY_FILE", "")
	viper.SetDefault("TWEET_TTL", 7)
	viper.AutomaticEnv()

	loglevel, err := log.ParseLevel(viper.GetString("LOGLEVEL"))
	if err != nil {
		loglevel = log.ErrorLevel
	}

	conf := Config{
		Host:             viper.GetString("HOST"),
		Port:             viper.GetInt("PORT"),
		Debug:            viper.GetInt("DEBUG"),
		Loglevel:         loglevel,
		DSN:              viper.GetString("DSN"),
		TwProxyPort:      viper.GetInt("TW_PROXY_PORT"),
		TwProxyHost:      viper.GetString("TW_PROXY_HOST"),
		Path:             viper.GetString("PATH"),
		Domain:           viper.GetString("DOMAIN"),
		MaxAge:           viper.GetInt("MAX_AGE"),
		Secure:           viper.GetBool("SECURE"),
		HttpOnly:         viper.GetBool("HTTP_ONLY"),
		AuthKey:          viper.GetString("auth-key"),
		EncryptKey:       viper.GetString("encrypt-key"),
		Testing:          viper.GetBool("TESTING"),
		TwConsumerKey:    viper.GetString("tw-consumer-key"),
		TwConsumerSecret: viper.GetString("tw-consumer-secret"),
		TwCallbackURL:    viper.GetString("TW_CALLBACK_URL"),
		TemplatePath:     viper.GetString("TEMPLATE_PATH"),
		MgDomain:         viper.GetString("MG_DOMAIN"),
		MgAPIKEY:         viper.GetString("MG_APIKEY"),
		From:             viper.GetString("FROM"),
		PemFile:          viper.GetString("PEM_FILE"),
		KeyFile:          viper.GetString("KEY_FILE"),
		TweetTTL:         viper.GetInt("TWEET_TTL"),
	}

	return conf
}
