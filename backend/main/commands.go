package main

import (
	"github.com/dmtr/mail_me_all/backend/app"
	log "github.com/sirupsen/logrus"
)

func CheckNewSubscriptions(a *app.App) {
	log.Info("Executing CheckNewSubscriptions command")

	err := a.UseCases.InitSubscriptions()
	if err != nil {
		log.Errorf("Got error executing command %s", err)
	}

	log.Info("Command CheckNewSubscriptions finished")
}
