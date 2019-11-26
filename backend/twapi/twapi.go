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
	page             = 1
	count            = 10
)

// UserInfo represents twitter user information
type UserInfo struct {
	TwitterID     string
	Name          string
	Email         string
	ScreenName    string
	ProfileIMGURL string
}

type Tweet struct {
	IDStr                string
	Text                 string
	FullText             string
	InReplyToStatusIDStr string
	InReplyToUserIDStr   string
	UserID               string
	UserName             string
	UserScreenName       string
	UserProfileImageUrl  string
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

func (t Twitter) SearchUsers(accessToken, accessSecret, twitterID, query string) ([]UserInfo, error) {
	client := t.getSession(accessToken, accessSecret, twitterID)
	includeEntities := false
	users, _, err := client.Users.Search(query, &tw.UserSearchParams{Page: page, Count: count, IncludeEntities: &includeEntities})

	if err != nil {
		log.Errorf("Got error calling twitter api: %s", err)
		return make([]UserInfo, 0, 0), err
	}

	res := make([]UserInfo, 0, len(users))
	for _, user := range users {
		u := UserInfo{
			TwitterID:     user.IDStr,
			Name:          user.Name,
			Email:         user.Email,
			ScreenName:    user.ScreenName,
			ProfileIMGURL: user.ProfileImageURLHttps,
		}
		res = append(res, u)
	}

	return res, err
}

func (t Twitter) GetUserTimeline(accessToken, accessSecret, twitterID, screenName string, sinceID int64, count int64) ([]Tweet, error) {
	client := t.getSession(accessToken, accessSecret, twitterID)

	trim := true
	if count == 0 {
		trim = false
	}

	params := tw.UserTimelineParams{
		ScreenName: screenName,
		TrimUser:   &trim,
		TweetMode:  "extended",
	}

	if sinceID != 0 {
		params.SinceID = sinceID
	}

	if count != 0 {
		params.Count = int(count)
	}

	tweets, _, err := client.Timelines.UserTimeline(&params)

	if err != nil {
		log.Errorf("Got error calling twitter api: %s", err)
		return []Tweet{}, err
	}

	res := make([]Tweet, 0, len(tweets))
	for _, tweet := range tweets {
		t := Tweet{
			IDStr:                tweet.IDStr,
			Text:                 tweet.Text,
			FullText:             tweet.FullText,
			InReplyToStatusIDStr: tweet.InReplyToStatusIDStr,
			InReplyToUserIDStr:   tweet.InReplyToUserIDStr,
			UserID:               tweet.User.IDStr,
			UserName:             tweet.User.Name,
			UserScreenName:       tweet.User.ScreenName,
			UserProfileImageUrl:  tweet.User.ProfileImageURL,
		}
		res = append(res, t)
	}

	return res, err
}
