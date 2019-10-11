package usecases

import (
	"context"

	"github.com/dmtr/mail_me_all/backend/db"
	"github.com/dmtr/mail_me_all/backend/models"
	log "github.com/sirupsen/logrus"

	pb "github.com/dmtr/mail_me_all/backend/rpc"
)

const (
	userCreationError   string = "Can not create user"
	tokenInsertionError string = "Can not save token"
)

type UserUseCase struct {
	UserDatastore models.UserDatastore
	RpcClient     pb.FbProxyServiceClient
}

func NewUserUseCase(datastore models.UserDatastore, client pb.FbProxyServiceClient) *UserUseCase {
	return &UserUseCase{UserDatastore: datastore, RpcClient: client}
}

func (u UserUseCase) SignInFB(ctx context.Context, userID string, accessToken string) error {
	log.Debugf("Sign in user %s", userID)

	newUser := pb.UserToken{UserId: userID, AccessToken: accessToken}

	confirmedUser, err := u.RpcClient.GetAccessToken(context.Background(), &newUser)
	if err != nil {
		log.Errorf("Can not get user access token, got error: %s", err)
		return NewUseCaseError(userCreationError)
	}

	if userID != confirmedUser.UserId {
		log.Warningf("Users ids do not match %s %s", confirmedUser.UserId, userID)
		return NewUseCaseError(userCreationError)
	}

	userInfo, err := u.RpcClient.GetUserInfo(context.Background(), confirmedUser)
	if err != nil {
		log.Errorf("Error %s", err)
	}
	user := models.User{
		Name:  userInfo.Name,
		FbID:  confirmedUser.UserId,
		Email: userInfo.Email,
	}

	if newUser, err := u.UserDatastore.InsertUser(ctx, user); err != nil {
		e, ok := err.(*db.DbError)
		if !ok {
			log.Errorf("Can not convert error to DbError: %s", err)
			return NewUseCaseError(userCreationError)
		}

		if e.PqError.Code != db.UniqueViolationErr {
			return NewUseCaseError(userCreationError)
		}
	} else {
		log.Debugf("New user %s", newUser)

		t := models.Token{
			UserID:  newUser.ID,
			FbToken: confirmedUser.AccessToken,
		}
		t.ExpiresAt = t.CalculateExpiresAt(confirmedUser.ExpiresIn)

		if token, err := u.UserDatastore.InsertToken(ctx, t); err != nil {
			return NewUseCaseError(tokenInsertionError)
		} else {
			log.Debugf("Token %s", token)
		}
	}
	return nil
}
