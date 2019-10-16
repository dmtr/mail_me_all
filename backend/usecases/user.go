package usecases

import (
	"context"

	"github.com/dmtr/mail_me_all/backend/db"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/google/uuid"
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

func (u UserUseCase) SignInFB(ctx context.Context, userID string, accessToken string) (models.User, error) {
	log.Debugf("Sign in user %s", userID)

	newUser := pb.UserToken{UserId: userID, AccessToken: accessToken}

	longToken, err := u.RpcClient.GetAccessToken(context.Background(), &newUser)
	if err != nil {
		log.Errorf("Can not get user access token, got error: %s", err)
		return models.User{}, NewUseCaseError(userCreationError)
	}

	if userID != longToken.UserId {
		log.Warningf("Users ids do not match %s %s", longToken.UserId, userID)
		return models.User{}, NewUseCaseError(userCreationError)
	}

	userInfo, err := u.RpcClient.GetUserInfo(context.Background(), longToken)
	if err != nil {
		log.Errorf("Error %s", err)
	}
	user := models.User{
		Name:  userInfo.Name,
		FbID:  longToken.UserId,
		Email: userInfo.Email,
	}

	if user, err = u.UserDatastore.InsertUser(ctx, user); err != nil {
		e, ok := err.(*db.DbError)
		if !ok {
			log.Errorf("Can not convert error to DbError: %s", err)
			return models.User{}, NewUseCaseError(userCreationError)
		}

		if e.PqError.Code != db.UniqueViolationErr {
			return models.User{}, NewUseCaseError(userCreationError)
		}
	} else {
		log.Debugf("New user %s", user)

		t := models.Token{
			UserID:  user.ID,
			FbToken: longToken.AccessToken,
		}
		t.ExpiresAt = t.CalculateExpiresAt(longToken.ExpiresIn)

		if token, err := u.UserDatastore.InsertToken(ctx, t); err != nil {
			return models.User{}, NewUseCaseError(tokenInsertionError)
		} else {
			log.Debugf("Token %s", token)
		}
	}
	return user, nil
}

func (u UserUseCase) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	return u.UserDatastore.GetUserByID(ctx, userID)
}
