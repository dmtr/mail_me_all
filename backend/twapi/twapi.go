package twapi

import (
	"sync"
	"time"

	tw "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	log "github.com/sirupsen/logrus"
)

const (
	sessionExpiresIn = 60 * 10
)

// UserInfo represents twitter user information
type UserInfo struct {
	TwitterID     string
	Name          string
	Email         string
	ScreenName    string
	ProfileIMGURL string
}

type Twitter struct {
	oauth1Config *oauth1.Config
	sessions     map[string]*tw.Client
	mux          *sync.Mutex
}

func NewTwitter(consumerKey, consumerSecret string) Twitter {
	return Twitter{
		oauth1Config: oauth1.NewConfig(consumerKey, consumerSecret),
		sessions:     make(map[string]*tw.Client),
		mux:          &sync.Mutex{},
	}
}

func (t Twitter) addSessionWithExpiration(twitterID string, client *tw.Client) {
	timer := time.NewTimer(sessionExpiresIn * time.Second)
	t.sessions[twitterID] = client
	go func() {
		<-timer.C
		t.deleteSession(twitterID)
	}()
}

func (t Twitter) deleteSession(twitterID string) {
	t.mux.Lock()
	defer t.mux.Unlock()
	log.Debugf("Deleting Twitter session for user %s", twitterID)
	delete(t.sessions, twitterID)
}

func (t Twitter) getClient(accessToken, accessSecret, twitterID string) *tw.Client {
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := t.oauth1Config.Client(oauth1.NoContext, token)
	return tw.NewClient(httpClient)
}

func (t Twitter) addSession(accessToken, accessSecret, twitterID string) *tw.Client {
	t.mux.Lock()
	defer t.mux.Unlock()
	client := t.getClient(accessToken, accessSecret, twitterID)
	t.addSessionWithExpiration(twitterID, client)
	return client
}

func (t Twitter) getSession(accessToken, accessSecret, twitterID string) *tw.Client {
	t.mux.Lock()
	defer t.mux.Unlock()
	client, ok := t.sessions[twitterID]
	if ok {
		return client
	}

	client = t.getClient(accessToken, accessSecret, twitterID)
	t.addSessionWithExpiration(twitterID, client)
	return client
}

func (t Twitter) GetUserInfo(accessToken, accessSecret, twitterID, screenName string) (UserInfo, error) {
	client := t.getSession(accessToken, accessSecret, twitterID)
	user, _, err := client.Users.Show(&tw.UserShowParams{
		ScreenName: screenName,
	})
	if err != nil {
		log.Errorf("Got error calling twitter api: %s", err)
		return UserInfo{}, err
	}

	u := UserInfo{
		TwitterID:     user.IDStr,
		Name:          user.Name,
		Email:         user.Email,
		ScreenName:    user.ScreenName,
		ProfileIMGURL: user.ProfileImageURLHttps,
	}
	return u, err
}
