package main

import (
	"github.com/dmtr/mail_me_all/backend/app"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const (
	initKey    = 1
	prepareKey = 2
	sendKey    = 3
)

func checkNewSubscriptions(a *app.App, ids ...uuid.UUID) {
	log.Info("Executing checkNewSubscriptions command")

	err := a.UseCases.InitSubscriptions(ids...)
	if err != nil {
		log.Errorf("Got error executing command %s", err)
	}

	log.Info("Command CheckNewSubscriptions finished")
}

func prepareSubscriptions(a *app.App, ids ...uuid.UUID) {
	log.Info("Executing prepareSubscriptions command")

	err := a.UseCases.PrepareSubscriptions(ids...)
	if err != nil {
		log.Errorf("Got error executing command %s", err)
	}

	log.Info("Command prepareSubscriptions finished")
}

func sendSubscriptions(a *app.App, ids ...uuid.UUID) {
	log.Info("Executing sendSubscriptions command")

	err := a.UseCases.SendSubscriptions(a.Conf.TemplatePath, ids...)
	if err != nil {
		log.Errorf("Got error executing command %s", err)
	}

	log.Info("Command sendSubscriptions finished")
}
