package main

import (
	"github.com/dmtr/mail_me_all/backend/fbwrapper"
	log "github.com/sirupsen/logrus"
)

func VerifyFbLogin(accessToken string, app *App) {
	log.Debugf("Verifying token %s", accessToken)
	fbwrapper.VerifyFbLogin(accessToken, app.Conf.FbAppID, app.Conf.AppSecret, app.Conf.FbRedirectURI)
}
