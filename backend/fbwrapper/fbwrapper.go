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

//AccessTokenResponse - long lived token response
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint   `json:"expires_in"`
}

//Facebook - facebook client
type Facebook struct {
	AppSecret string
	AppID     string
	App       *fb.App
}

//NewFacebook - returns new facebook client
func NewFacebook(appID string, appSecret string, redirectURI string) Facebook {
	var app = fb.New(appID, appSecret)
	app.RedirectUri = redirectURI

	return Facebook{
		AppSecret: appSecret,
		AppID:     appID,
		App:       app,
	}
}

//VerifyFbToken - check if access token is valid
func (f Facebook) VerifyFbToken(accessToken string) (userid string, err error) {
	session := f.App.Session(accessToken)
	return session.User()
}

//GenerateLongLivedToken - generates long lived token
func (f Facebook) GenerateLongLivedToken(accessToken string) (AccessTokenResponse, error) {
	var response AccessTokenResponse

	req, err := http.NewRequest("GET", fbTokenUrl, nil)
	if err != nil {
		log.Errorf("%s", err)
		return response, err
	}

	q := req.URL.Query()
	q.Add("grant_type", grantType)
	q.Add("client_id", f.AppID)
	q.Add("client_secret", f.AppSecret)
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
