package usecases

import (
	"github.com/dmtr/mail_me_all/backend/models"
	pb "github.com/dmtr/mail_me_all/backend/rpc"
)

// SystemUseCase implementation
type SystemUseCase struct {
	UserDatastore models.UserDatastore
	RpcClient     pb.TwProxyServiceClient
}

func NewSystemUseCase(datastore models.UserDatastore, client pb.TwProxyServiceClient) *SystemUseCase {
	return &SystemUseCase{UserDatastore: datastore, RpcClient: client}
}

func (s SystemUseCase) InitSubscriptions() error {
	return nil
}
