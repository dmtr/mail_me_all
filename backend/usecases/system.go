package usecases

import (
	"context"
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

func (s SystemUseCase) initSubscription(subscriptionID uuid.UUID, wg *sync.WaitGroup) {
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

		log.Infof("tweets: %v", tweets)
	}

}

func (s SystemUseCase) InitSubscriptions(ids ...uuid.UUID) error {
	var subscriptions []uuid.UUID
	var err error

	if len(ids) == 0 {
		subscriptions, err = s.UserDatastore.GetNewSubscriptionsIDs(context.Background())
		if err != nil {
			return err
		}
	} else {
		subscriptions = ids
	}

	log.Infof("Got subscriptions %s", subscriptions)

	var wg sync.WaitGroup
	for _, id := range subscriptions {
		wg.Add(1)
		go s.initSubscription(id, &wg)
	}

	wg.Wait()
	return err
}
