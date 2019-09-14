package main

import (
	"os"
	"time"

	"github.com/dmtr/mail_me_all/backend/api"
	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/db"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/dmtr/mail_me_all/backend/usecases"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	retry time.Duration = 20
)

// App represents application
type App struct {
	Router *gin.Engine
	Conf   *config.Config
	Db     *sqlx.DB
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
func GetApp(withAPI bool, withDB bool) *App {
	log.Infoln("Loading Config")
	conf := config.GetConfig()
	initLogger(conf.Loglevel)
	log.Infof("Config loaded %v", conf)

	var db_ *sqlx.DB

	if withDB {
		var err error
		db_, err = db.ConnectDb(conf.DSN, retry*time.Second)
		if err != nil {
			log.Fatalf("Can't connect to database %s", err)
			os.Exit(1)
		}
	}

	fn := func() {
		log.Info("Closing.")
		if withDB {
			db_.Close()
		}
	}

	userDatastore := db.NewUserDatastore(db_)
	userUseCase := usecases.NewUserUseCase(userDatastore)
	usecases := models.NewUseCases(userUseCase)

	var router *gin.Engine
	if withAPI {

		if conf.Debug == 0 {
			log.Info("Release mode")
			gin.SetMode(gin.ReleaseMode)
		}
		router = api.GetRouter(usecases)
	}

	return &App{
		Router: router,
		Conf:   &conf,
		Db:     db_,
		Close:  fn,
	}
}
