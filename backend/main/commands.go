package main

import (
	"github.com/dmtr/mail_me_all/backend/app"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CheckNewSubscriptions(a *app.App, ids ...uuid.UUID) {
	log.Info("Executing CheckNewSubscriptions command")

	err := a.UseCases.InitSubscriptions(ids...)
	if err != nil {
		log.Errorf("Got error executing command %s", err)
	}

	log.Info("Command CheckNewSubscriptions finished")
}
