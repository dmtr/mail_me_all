package app

import (
	"os"
	"time"

	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/db"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	retry time.Duration = 20
)

// App represents application
type App struct {
	Router *gin.Engine
	Conf   *config.Config
	Close  func()
}

func initLogger(loglevel log.Level) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true})
	log.SetOutput(os.Stdout)
	if loglevel == 0 {
		loglevel = log.ErrorLevel
	}
	log.SetLevel(loglevel)
	log.SetReportCaller(true)
}

// GetApp - returns app
func GetApp() App {
	log.Infoln("Loading Config")
	conf := config.GetConfig()
	initLogger(conf.Loglevel)
	log.Infof("Config loaded %v", conf)

	if conf.Debug == 0 {
		log.Info("Release mode")
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := db.ConnectDb(conf.DSN, retry*time.Second)
	if err != nil {
		log.Fatalf("Can't connect to database %s", err)
		os.Exit(1)
	}

	fn := func() { log.Info("Closing."); db.Close() }
	return App{
		Router: GetRouter(db),
		Conf:   &conf,
		Close:  fn,
	}
}
