package main

import (
	"github.com/dmtr/mail_me_all/backend/app"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

	err := a.UseCases.SendSubscriptions(ids...)
	if err != nil {
		log.Errorf("Got error executing command %s", err)
	}

	log.Info("Command sendSubscriptions finished")
}

func sendConfirmationEmail(a *app.App) {
	log.Info("Executing sendConfirmationEmail command")

	err := a.UseCases.SendConfirmationEmail()
	if err != nil {
		log.Errorf("Got error executing command %s", err)
	}

	log.Info("Command sendConfirmationEmail finished")
}

func removeOldTweets(a *app.App) {
	log.Info("Executing removeOldTweets command")

	err := a.UseCases.RemoveOldTweets()
	if err != nil {
		log.Errorf("Got error executing command %s", err)
	}

	log.Info("Command removeOldTweets finished")
}
