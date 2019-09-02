package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/dmtr/mail_me_all/backend/app"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := app.GetApp()
	defer app.Close()

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
