package usecases

import "github.com/dmtr/mail_me_all/backend/models"
import log "github.com/sirupsen/logrus"
import pb "github.com/dmtr/mail_me_all/backend/rpc"

const (
	userCreationError string = "Can not create user"
)

type UserUseCase struct {
	UserDatastore models.UserDatastore
	RpcClient     pb.FbProxyServiceClient
}

func NewUserUseCase(datastore models.UserDatastore, client pb.FbProxyServiceClient) *UserUseCase {
	return &UserUseCase{UserDatastore: datastore, RpcClient: client}
}

func (u UserUseCase) SignInFB(userID string, accessToken string) error {
	log.Debugf("Sign in user %s", userID)
	//if err := u.UserDatastore.CreateUser(user); err != nil {
	//	return NewUseCaseError(userCreationError)
	//}
	return nil
}
