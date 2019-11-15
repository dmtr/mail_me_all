package app

import (
	"os"
	"time"

	"github.com/dmtr/mail_me_all/backend/api"
	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/db"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/dmtr/mail_me_all/backend/rpc"
	"github.com/dmtr/mail_me_all/backend/usecases"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/dmtr/mail_me_all/backend/rpc"
)

const (
	retry time.Duration = 20
)

// App represents application
type App struct {
	Router   *gin.Engine
	Conf     *config.Config
	Db       *sqlx.DB
	UseCases *models.UseCases
	Close    func()
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

func getUseCases(db_ *sqlx.DB, client pb.TwProxyServiceClient) *models.UseCases {
	userDatastore := db.NewUserDatastore(db_)
	userUseCase := usecases.NewUserUseCase(userDatastore, client)
	systemUseCase := usecases.NewSystemUseCase(userDatastore, client)
	return models.NewUseCases(userUseCase, systemUseCase)
}

// GetApp - returns app
func GetApp(withAPI bool, opts ...bool) *App {
	log.Infoln("Loading Config")
	conf := config.GetConfig()
	initLogger(conf.Loglevel)
	log.Infof("Config loaded %v", conf)

	withDB := false
	withRpcConn := false
	withUseCases := false

	if withAPI {
		withDB = true
		withRpcConn = true
		withUseCases = true
	} else {
		for i, v := range opts {
			if i == 0 {
				withDB = v
			} else if i == 1 {
				withRpcConn = v
			} else if i == 2 {
				withUseCases = v
			}
		}
	}

	var db_ *sqlx.DB

	if withDB {
		var err error
		db_, err = db.ConnectDb(conf.DSN, retry*time.Second)
		if err != nil {
			log.Fatalf("Can't connect to database %s", err)
			os.Exit(1)
		}
	}

	var conn *grpc.ClientConn

	if withRpcConn {
		var err error
		conn, err = rpc.GetRpcConection(&conf)
		if err != nil {
			log.Fatalf("Can't connect to rpc sever %s", err)
			os.Exit(1)
		}

	}

	fn := func() {
		log.Info("Closing.")
		if withDB {
			db_.Close()
		}
		if withRpcConn {
			conn.Close()
		}
	}

	var usecases *models.UseCases

	if withUseCases {
		client := rpc.GetRpcClient(conn)
		usecases = getUseCases(db_, client)
	}

	var router *gin.Engine

	if withAPI {
		if conf.Debug == 0 {
			log.Info("Release mode")
			gin.SetMode(gin.ReleaseMode)
		}
		router = api.GetRouter(&conf, db_, usecases)
	}

	return &App{
		Router:   router,
		Conf:     &conf,
		Db:       db_,
		UseCases: usecases,
		Close:    fn,
	}
}
