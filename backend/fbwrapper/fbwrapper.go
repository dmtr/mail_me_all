package fbwrapper

import (
	"encoding/json"
	"net/http"
	"time"

	fb "github.com/huandu/facebook"
	log "github.com/sirupsen/logrus"
)

const (
	fbTokenUrl = "https://graph.facebook.com/v4.0/oauth/access_token"
	timeout    = time.Duration(10 * time.Second)
	grantType  = "fb_exchange_token"
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

func GenerateLongLivedToken(accessToken string, appID string, appSecret string) (string, error) {
	req, err := http.NewRequest("GET", fbTokenUrl, nil)
	if err != nil {
		log.Errorf("%s", err)
		return "", err
	}

	q := req.URL.Query()
	q.Add("grant_type", grantType)
	q.Add("client_id", appID)
	q.Add("client_secret", appSecret)
	q.Add("fb_exchange_token", accessToken)
	req.URL.RawQuery = q.Encode()

	client := http.Client{
		Timeout: timeout,
	}
	log.Debugf("Url: %s", req.URL.String())
	resp, err := client.Do(req)
	log.Debugf("Response code %s", resp.StatusCode)
	if err != nil {
		log.Errorf("%s", err)
		return "", err

	}

	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	log.Debugf("result %v", result)
	return result["access_token"].(string), nil

}
