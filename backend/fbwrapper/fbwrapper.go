package fbwrapper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
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

//UserInfo - user data
type UserInfo struct {
	UserID    string
	FirstName string
	Email     string
}

//Facebook - facebook client
type Facebook struct {
	AppSecret string
	AppID     string
	App       *fb.App
	sessions  map[string]*fb.Session
	mux       *sync.Mutex
}

//NewFacebook - returns new facebook client
func NewFacebook(appID string, appSecret string, redirectURI string) Facebook {
	var app = fb.New(appID, appSecret)
	app.RedirectUri = redirectURI

	return Facebook{
		AppSecret: appSecret,
		AppID:     appID,
		App:       app,
		sessions:  make(map[string]*fb.Session),
		mux:       &sync.Mutex{},
	}
}

func (f Facebook) addSession(accessToken string) *fb.Session {
	f.mux.Lock()
	defer f.mux.Unlock()

	s := f.App.Session(accessToken)
	userID, err := s.User()
	if err != nil {
		return s
	}

	userSession, ok := f.sessions[userID]
	if ok {
		return userSession
	}

	err = s.Validate()
	if err == nil {
		f.sessions[userID] = s
	}
	return s
}

func (f Facebook) getSession(userID string, accessToken string) *fb.Session {
	f.mux.Lock()
	defer f.mux.Unlock()
	userSession, ok := f.sessions[userID]
	if ok {
		return userSession
	}
	userSession = f.App.Session(accessToken)
	err := userSession.Validate()
	if err == nil {
		f.sessions[userID] = userSession
	}
	return userSession
}

func (f Facebook) deleteSession(userID string) {
	f.mux.Lock()
	defer f.mux.Unlock()
	delete(f.sessions, userID)
}

//VerifyFbToken - check if access token is valid
func (f Facebook) VerifyFbToken(accessToken string) (userid string, err error) {
	s := f.addSession(accessToken)
	return s.User()
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

//GetUserInfo - returns user info
func (f Facebook) GetUserInfo(userID string, accessToken string) (UserInfo, error) {
	var u UserInfo
	s := f.getSession(userID, accessToken)
	res, err := s.Get(fmt.Sprintf("/%s", userID), fb.Params{
		"fields":       "first_name,email",
		"access_token": accessToken,
	})

	if err != nil {
		log.Errorf("Cant get user %s info, got error %s", userID, err)
	} else {
		res.Decode(&u)
	}

	u.UserID = userID
	log.Debugf("User info %v", u)
	return u, err
}
