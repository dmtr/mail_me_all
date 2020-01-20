package usecases

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/models"
	pb "github.com/dmtr/mail_me_all/backend/rpc"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

const (
	initKey    = 1
	prepareKey = 2
	sendKey    = 3
)

var once sync.Once
var shortenerRegex *regexp.Regexp

func getShortenerRegexp() *regexp.Regexp {
	once.Do(func() {
		r, err := regexp.Compile("(https://t.co/[A-Za-z0-9]+)")
		if err != nil {
			log.Errorf("Can't compile regex %s", err)
		} else {
			shortenerRegex = r
		}
	})
	return shortenerRegex
}

type JWTClaims struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// SystemUseCase implementation
type SystemUseCase struct {
	UserDatastore models.UserDatastore
	RpcClient     pb.TwProxyServiceClient
	Conf          *config.Config
	EmailSender   models.EmailSender
}

// NewSystemUseCase creates a new SystemUseCase
func NewSystemUseCase(datastore models.UserDatastore, client pb.TwProxyServiceClient, conf *config.Config, emailSender models.EmailSender) *SystemUseCase {
	return &SystemUseCase{
		UserDatastore: datastore,
		RpcClient:     client,
		Conf:          conf,
		EmailSender:   emailSender}
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func merge(channels []<-chan models.Tweet) <-chan models.Tweet {
	var wg sync.WaitGroup
	out := make(chan models.Tweet)

	output := func(c <-chan models.Tweet) {
		for t := range c {
			out <- t
		}
		wg.Done()
	}

	wg.Add(len(channels))

	for _, c := range channels {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (s SystemUseCase) initSubscription(subscriptionID uuid.UUID, users []string, wg *sync.WaitGroup) {
	defer wg.Done()
	subscription, err := s.UserDatastore.GetSubscription(context.Background(), subscriptionID)
	if err != nil {
		log.Errorf("Can not get subscription %s, got error %s", subscriptionID, err)
		return
	}
	log.Infof("Got subscription %s", subscription)

	user, err := s.UserDatastore.GetTwitterUser(context.Background(), subscription.UserID)
	if err != nil {
		log.Errorf("Can not get user %s, got error %s", subscription.UserID, err)
		return
	}

	log.Infof("Got user %s", user)

	for _, u := range subscription.UserList {

		_, found := find(users, u.TwitterID)
		if !found {
			continue
		}

		req := pb.UserTimelineRequest{
			AccessToken:  user.AccessToken,
			AccessSecret: user.TokenSecret,
			TwitterId:    user.TwitterID,
			ScreenName:   u.ScreenName,
			SinceId:      0,
			Count:        1}

		tweets, err := s.RpcClient.GetUserTimeline(context.Background(), &req)
		if err != nil {
			log.Errorf("Can not get timeline for user %s", u)
		}

		log.Debugf("tweets: %v", tweets)

		err = s.UserDatastore.InsertSubscriptionUserState(context.Background(), subscriptionID, u.TwitterID, tweets.Tweets[0].IdStr)
		if err != nil {
			log.Errorf("Can not insert subscription_user_state, got error %s", err)
		}
	}
}

func (s SystemUseCase) InitSubscriptions(ids ...uuid.UUID) error {
	lock, err := s.UserDatastore.AcquireLock(context.Background(), initKey)
	if err != nil {
		return err
	}

	if !lock {
		return errors.New("Can not acquire lock")
	}

	defer func() {
		_, err = s.UserDatastore.ReleaseLock(context.Background(), initKey)
		if err != nil {
			log.Errorf("Can not release lock %s", err)
		}
	}()

	subscriptions, err := s.UserDatastore.GetNewSubscriptionsUsers(context.Background(), ids...)
	if err != nil {
		return err
	}

	log.Infof("Got subscriptions %s", subscriptions)

	var wg sync.WaitGroup
	for subscription, users := range subscriptions {
		wg.Add(1)
		go s.initSubscription(subscription, users, &wg)
	}

	wg.Wait()
	return err
}

func (s SystemUseCase) PrepareSubscriptions(ids ...uuid.UUID) error {
	lock, err := s.UserDatastore.AcquireLock(context.Background(), prepareKey)
	if err != nil {
		return err
	}

	if !lock {
		return errors.New("Can not acquire lock")
	}

	defer func() {
		_, err = s.UserDatastore.ReleaseLock(context.Background(), prepareKey)
		if err != nil {
			log.Errorf("Can not release lock %s", err)
		}
	}()

	var subscriptions []uuid.UUID

	if len(ids) == 0 {
		subscriptions, err = s.UserDatastore.GetTodaySubscriptionsIDs(context.Background())
		if err != nil {
			return err
		}
	} else {
		subscriptions = ids
	}

	log.Infof("Got subscriptions %s", subscriptions)

	var wg sync.WaitGroup
	for _, id := range subscriptions {
		state, err := s.UserDatastore.InsertSubscriptionState(
			context.Background(), models.SubscriptionState{SubscriptionID: id, Status: models.Preparing})
		if err != nil {
			log.Errorf("Can not insert subscription state got error %s", err)
			continue
		}

		subscription, err := s.UserDatastore.GetSubscription(context.Background(), id)
		if err != nil {
			log.Errorf("Can not get subscription %s, got error %s", id, err)
			continue
		}
		log.Infof("Got subscription %s", subscription)

		user, err := s.UserDatastore.GetTwitterUser(context.Background(), subscription.UserID)
		if err != nil {
			log.Errorf("Can not get user %s, got error %s", subscription.UserID, err)
			continue
		}

		log.Infof("Got user %s", user)

		wg.Add(1)
		go s.prepareSubscription(subscription, user, state, &wg)
	}

	wg.Wait()

	return err
}

func (s SystemUseCase) prepareSubscription(subscription models.Subscription, user models.TwitterUser, subscriptionState models.SubscriptionState, wg *sync.WaitGroup) {
	defer wg.Done()

	subscriptionUserTweets, err := s.UserDatastore.GetSubscriptionUserTweets(context.Background(), subscription.ID)
	if err != nil {
		log.Errorf("Can't get subscription user' tweets %s", err)
		return
	}

	channels := make([]<-chan models.Tweet, 0)
	for _, u := range subscription.UserList {
		ch := s.getTweets(subscriptionUserTweets, u, user.AccessToken, user.TokenSecret, user.TwitterID, subscription.IgnoreRT, subscription.IgnoreReplies)
		channels = append(channels, ch)
	}

	for t := range merge(channels) {
		log.Infof("Got tweet %s", t.Tweet.FullText)
		_, err := s.UserDatastore.InsertTweet(context.Background(), t, subscriptionState.ID)
		if err != nil {
			log.Errorf("Can't insert tweet %s", err)
		}
	}

	subscriptionState.Status = models.Ready
	_, err = s.UserDatastore.UpdateSubscriptionState(context.Background(), subscriptionState)
	if err != nil {
		log.Errorf("Can't update subscription state %s  %s", subscriptionState.String(), err)
	}
}

func (s SystemUseCase) getTweets(subscriptionUserTweets models.SubscriptionUserTweets, user models.TwitterUserSearchResult, accessToken, tokenSecret, twitterID string, ignoreRT, ignoreReplies bool) <-chan models.Tweet {
	ch := make(chan models.Tweet)

	lastTweet, ok := subscriptionUserTweets.Tweets[user.TwitterID]
	if !ok {
		log.Errorf("Can't find last tweet for user %s", user)
		close(ch)
		return ch
	}

	sinceID, err := strconv.ParseInt(lastTweet.LastTweetID, 10, 64)
	if err != nil {
		log.Errorf("Can't parse last tweet id %s %s", lastTweet, err)
		close(ch)
		return ch
	}

	go func() {
		req := pb.UserTimelineRequest{
			AccessToken:   accessToken,
			AccessSecret:  tokenSecret,
			TwitterId:     twitterID,
			ScreenName:    user.ScreenName,
			SinceId:       sinceID,
			IgnoreRt:      ignoreRT,
			IgnoreReplies: ignoreReplies,
		}

		tweets, err := s.RpcClient.GetUserTimeline(context.Background(), &req)
		if err != nil {
			log.Errorf("Can not get timeline for user %s", user)
		}

		for _, t := range tweets.Tweets {
			ch <- models.Tweet{
				TweetID: t.IdStr,
				Tweet: models.TweetAttrs{
					IdStr:                t.IdStr,
					Text:                 t.Text,
					FullText:             t.FullText,
					InReplyToStatusIdStr: t.InReplyToStatusIdStr,
					InReplyToUserIdStr:   t.InReplyToUserIdStr,
					UserId:               t.UserId,
					UserName:             t.UserName,
					UserScreenName:       t.UserScreenName,
					UserProfileImageUrl:  t.UserProfileImageUrl,
				},
			}
		}
		close(ch)
	}()

	return ch
}

func (s SystemUseCase) SendSubscriptions(ids ...uuid.UUID) error {
	lock, err := s.UserDatastore.AcquireLock(context.Background(), sendKey)
	if err != nil {
		return err
	}

	if !lock {
		return errors.New("Can not acquire lock")
	}

	defer func() {
		_, err = s.UserDatastore.ReleaseLock(context.Background(), sendKey)
		if err != nil {
			log.Errorf("Can not release lock %s", err)
		}
	}()

	states, err := s.UserDatastore.GetReadySubscriptionsStates(context.Background(), ids...)
	if err != nil {
		return err
	}

	log.Infof("Got subscriptions %v", states)

	r := getShortenerRegexp()
	shortener := func(s string) template.HTML {
		return template.HTML(r.ReplaceAllStringFunc(s, func(t string) string { return fmt.Sprintf("<a href=\"%s\">%s</a>", t, t) }))
	}

	funcMap := template.FuncMap{
		"shortener": shortener,
	}

	templatePath := s.Conf.TemplatePath
	tmpl := template.Must(template.New("mail.html").Funcs(funcMap).ParseFiles(filepath.Join(templatePath, "mail.html")))

	var wg sync.WaitGroup

	for _, st := range states {
		st.Status = models.Sending
		state, err := s.UserDatastore.UpdateSubscriptionState(context.Background(), st)
		if err != nil {
			log.Errorf("Can not update subscription state got error %s", err)
			continue
		}

		subscription, err := s.UserDatastore.GetSubscription(context.Background(), state.SubscriptionID)
		if err != nil {
			log.Errorf("Can not get subscription %s, got error %s", state.SubscriptionID, err)
			continue
		}
		log.Infof("Got subscription %s", subscription)

		userEmail := models.UserEmail{
			UserID: subscription.UserID,
			Email:  subscription.Email,
		}

		email, err := s.UserDatastore.GetUserEmail(context.Background(), userEmail)
		if err != nil {
			log.Errorf("Can not get UserEmail %s, got error %s", userEmail, err)
			continue
		}

		if email.Status != models.EmailStatusConfirmed {
			log.Errorf("Email is not confirmed %s", email)
			continue
		}

		user, err := s.UserDatastore.GetTwitterUser(context.Background(), subscription.UserID)
		if err != nil {
			log.Errorf("Can not get user %s, got error %s", subscription.UserID, err)
			continue
		}

		log.Infof("Got user %s", user)

		wg.Add(1)
		go s.sendSubscription(subscription, state, tmpl, &wg)
	}

	wg.Wait()

	err = s.UserDatastore.UpdateSubscriptionUserStateTweets(context.Background())

	return err
}

func (s SystemUseCase) sendSubscription(subscription models.Subscription, subscriptionState models.SubscriptionState, tmpl *template.Template, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Infof("SubscriptionState %+v", subscriptionState)

	tweets, err := s.UserDatastore.GetSubscriptionTweets(context.Background(), subscriptionState.ID)
	if err != nil {
		log.Errorf("Can not get tweets for subscription %s, got error %s", subscription, err)
		subscriptionState.Status = models.Failed
		_, err = s.UserDatastore.UpdateSubscriptionState(context.Background(), subscriptionState)
		if err != nil {
			log.Errorf("Can not update subscription state got error %s", err)
		}
		return
	}

	if len(tweets) == 0 {
		log.Warnf("No tweets found for subscription %s", subscription)
		subscriptionState.Status = models.Sent
		_, err = s.UserDatastore.UpdateSubscriptionState(context.Background(), subscriptionState)
		if err != nil {
			log.Errorf("Can not update subscription state got error %s", err)
		}
		return
	}

	type TemplateData struct {
		Tweets []models.Tweet
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, TemplateData{Tweets: tweets})
	if err != nil {
		log.Errorf("err %s", err)
		subscriptionState.Status = models.Failed
		_, err = s.UserDatastore.UpdateSubscriptionState(context.Background(), subscriptionState)
		if err != nil {
			log.Errorf("Can not update subscription state got error %s", err)
		}
		return
	}

	log.Debugf("html %s", buf.String())

	err = s.EmailSender.Send(s.Conf.From, subscription.Email, subscription.GetSubject(), buf.String())

	if err != nil {
		subscriptionState.Status = models.Failed
	} else {
		subscriptionState.Status = models.Sent
	}
	_, err = s.UserDatastore.UpdateSubscriptionState(context.Background(), subscriptionState)
	if err != nil {
		log.Errorf("Can not update subscription state got error %s", err)
		return
	}
}

func (s SystemUseCase) GetToken(email, userID string) (string, error) {
	exp := time.Now().Unix() + 24*60*60
	claims := JWTClaims{
		email,
		userID,
		jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(s.Conf.EncryptKey))
	if err != nil {
		return "", err
	}

	return ss, err
}

func (s SystemUseCase) getEmailConfirmationLink(email, userID string) (string, error) {
	token, err := s.GetToken(email, userID)
	if err != nil {
		return "", err
	}

	link := &url.URL{
		Scheme:   "https",
		Host:     s.Conf.Domain,
		Path:     "confirm/email",
		RawQuery: fmt.Sprintf("token=%s", token),
	}

	return link.String(), err
}

func (s SystemUseCase) SendConfirmationEmail() error {
	emails, err := s.UserDatastore.GetUserEmails(context.Background(), models.EmailStatusNew)
	tmpl := template.Must(template.New("confirm.html").ParseFiles(filepath.Join(s.Conf.TemplatePath, "confirm.html")))

	for _, email := range emails {
		var buf strings.Builder
		link, err := s.getEmailConfirmationLink(email.Email, email.UserID.String())
		if err != nil {
			log.Errorf("Can not get confirmation link: %s", err)
		}

		type TemplateData struct {
			ConfirmationLink string
		}

		err = tmpl.Execute(&buf, TemplateData{ConfirmationLink: link})
		if err != nil {
			log.Errorf("Can not execute template: %s", err)
		}

		err = s.EmailSender.Send(s.Conf.From, email.Email, ConfirmationEmailSubj, buf.String())
		if err == nil {
			email.Status = models.EmailStatusSent
			_, err = s.UserDatastore.UpdateUserEmail(context.Background(), email)
			if err != nil {
				log.Errorf("Can not update user email: %s", err)
			}
		} else {
			log.Errorf("Can not send email: %s", err)
		}
	}
	return err
}
