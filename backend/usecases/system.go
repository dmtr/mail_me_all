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

func initSubscription(subscriptionID uuid.UUID, datastore models.UserDatastore, wg *sync.WaitGroup) {
	defer wg.Done()
	s, err := datastore.GetSubscription(context.Background(), subscriptionID)
	if err != nil {
		log.Errorf("Can not get subscription %s, got error %s", subscriptionID, err)
		return
	}
	log.Infof("Got subscription %s", s)
}

func (s SystemUseCase) InitSubscriptions() error {
	ctx := context.Background()
	subscriptions, err := s.UserDatastore.GetNewSubscriptionsIDs(ctx)

	if err != nil {
		return err
	}
	log.Infof("Got subscriptions %s", subscriptions)

	var wg sync.WaitGroup
	for _, id := range subscriptions {
		wg.Add(1)
		go initSubscription(id, s.UserDatastore, &wg)
	}

	wg.Wait()
	return err
}
