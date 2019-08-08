package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/routes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf := config.GetConfig()
	if conf.Debug == 0 {
		log.Info("Release mode")
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	routes.RegisterRoutes(router)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Handler: router,
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
