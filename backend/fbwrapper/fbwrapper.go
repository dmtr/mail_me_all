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

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint   `json:"expires_in"`
}

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

func GenerateLongLivedToken(accessToken string, appID string, appSecret string) (AccessTokenResponse, error) {
	var response AccessTokenResponse

	req, err := http.NewRequest("GET", fbTokenUrl, nil)
	if err != nil {
		log.Errorf("%s", err)
		return response, err
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

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Errorf("Response code %s, error %s", resp.Status, err)
		return response, err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&response); err != nil {
		log.Errorf("Got error decoding response: %s", err)
	}
	return response, nil
}
