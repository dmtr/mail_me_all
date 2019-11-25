package usecases

import (
	"context"
	"strconv"
	"sync"

	"github.com/dmtr/mail_me_all/backend/models"
	pb "github.com/dmtr/mail_me_all/backend/rpc"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

// SystemUseCase implementation
type SystemUseCase struct {
	UserDatastore models.UserDatastore
	RpcClient     pb.TwProxyServiceClient
}

func NewSystemUseCase(datastore models.UserDatastore, client pb.TwProxyServiceClient) *SystemUseCase {
	return &SystemUseCase{UserDatastore: datastore, RpcClient: client}
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
	var subscriptions []uuid.UUID
	var err error

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
		ch := s.getTweets(subscriptionUserTweets, u, user.AccessToken, user.TokenSecret, user.TwitterID)
		channels = append(channels, ch)
	}

	for t := range merge(channels) {
		log.Infof("Got tweet %s", t)
		_, err := s.UserDatastore.InsertTweet(context.Background(), t, subscriptionState.ID)
		if err != nil {
			log.Errorf("Can't insert tweet %s", err)
		}

		subscriptionState.Status = models.Ready
		_, err = s.UserDatastore.UpdateSubscriptionState(context.Background(), subscriptionState)
		if err != nil {
			log.Errorf("Can't update subscription state %s  %s", subscriptionState, err)
		}

	}

}

func (s SystemUseCase) getTweets(subscriptionUserTweets models.SubscriptionUserTweets, user models.TwitterUserSearchResult, accessToken, tokenSecret, twitterID string) <-chan models.Tweet {
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
			AccessToken:  accessToken,
			AccessSecret: tokenSecret,
			TwitterId:    twitterID,
			ScreenName:   user.ScreenName,
			SinceId:      sinceID,
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
	states, err := s.UserDatastore.GetReadySubscriptionsStates(context.Background(), ids...)
	if err != nil {
		return err
	}

	log.Infof("Got subscriptions %s", states)

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

		user, err := s.UserDatastore.GetTwitterUser(context.Background(), subscription.UserID)
		if err != nil {
			log.Errorf("Can not get user %s, got error %s", subscription.UserID, err)
			continue
		}

		log.Infof("Got user %s", user)

		wg.Add(1)
		go s.sendSubscription(subscription, state, &wg)
	}

	wg.Wait()
	return err
}

func (s SystemUseCase) sendSubscription(subscription models.Subscription, subscriptionState models.SubscriptionState, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Infof("SubscriptionState %+v", subscriptionState)

	tweets, err := s.UserDatastore.GetSubscriptionTweets(context.Background(), subscriptionState.ID)
	if err != nil {
		log.Errorf("Can not get tweets for subscription %s, got error %s", subscription, err)
	}

	for _, tweet := range tweets {
		log.Infof("Tweet %+v", tweet)
	}

}
