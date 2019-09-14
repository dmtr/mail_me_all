package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	runAPI        string = "api"
	verifyFbLogin string = "verify-fb-login"
)

func startApiServer(app *App) {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", app.Conf.Host, app.Conf.Port),
		Handler: app.Router,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Info("Recieve interrupt signal")
		err := server.Close()
		if err != nil {
			log.Errorf("Web server closed : %v", err)
		}

	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Info("Web server shutdown complete")
		} else {
			log.Errorf("Web server closed unexpect: %s", err)
		}
	}
	log.Info("Exiting")
}

func main() {
	flag.String("app-secret", "", "app secret")
	var accessToken *string = flag.String("access-token", "", "access token")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

	cmd := flag.Arg(0)
	if cmd == "" {
		cmd = runAPI
	}

	var app *App
	if cmd == runAPI {
		app = GetApp(true, true)
	} else if cmd == verifyFbLogin {
		app = GetApp(false, false)
		VerifyFbLogin(*accessToken, app)
	} else {
		fmt.Printf("Unknown command %s", cmd)
		os.Exit(1)
	}

	if app != nil {
		defer app.Close()
	}
}
