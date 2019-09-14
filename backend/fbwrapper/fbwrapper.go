package fbwrapper

import (
	fb "github.com/huandu/facebook"
	log "github.com/sirupsen/logrus"
)

func VerifyFbLogin(accessToken string, appID string, appSecret string, FbRedirectURI string) {
	var globalApp = fb.New(appID, appSecret)
	globalApp.RedirectUri = FbRedirectURI
	session := globalApp.Session(accessToken)
	err := session.Validate()
	if err != nil {
		log.Errorf("Invalid access token, got error %s", err)
	} else {
		log.Info("Token is valid")
		user, _ := session.User()
		log.Infof("User is %v", user)
	}
}
