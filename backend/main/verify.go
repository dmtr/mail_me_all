package main

import (
	"github.com/dmtr/mail_me_all/backend/fbwrapper"
	log "github.com/sirupsen/logrus"
)

func VerifyFbLogin(accessToken string, f fbwrapper.Facebook) {
	user, err := f.VerifyFbToken(accessToken)
	if err != nil {
		log.Errorf("Invalid access token: error %s", err)
	} else {
		log.Infof("Token is valid, user id %s", user)
	}
}

func GenerateFbToken(accessToken string, f fbwrapper.Facebook) {
	token, err := f.GenerateLongLivedToken(accessToken)
	if err != nil {
		log.Errorf("Got error %s", err)
	} else {
		log.Debugf("Got token %v", token)
	}
}
