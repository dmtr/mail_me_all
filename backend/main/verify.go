package main

import (
	"github.com/dmtr/mail_me_all/backend/fbwrapper"
	log "github.com/sirupsen/logrus"
)

func VerifyFbLogin(accessToken string, app *App) {
	user, err := fbwrapper.VerifyFbToken(accessToken, app.Conf.FbAppID, app.Conf.AppSecret, app.Conf.FbRedirectURI)
	if err != nil {
		log.Errorf("Invalid access token: error %s", err)
	} else {
		log.Infof("Token is valid, user id %s", user)
	}
}

func GenerateFbToken(accessToken string, app *App) {
	token, err := fbwrapper.GenerateLongLivedToken(accessToken, app.Conf.FbAppID, app.Conf.AppSecret)
	if err != nil {
		log.Errorf("Got error %s", err)
	} else {
		log.Debugf("Got token %v", token)
	}
}
