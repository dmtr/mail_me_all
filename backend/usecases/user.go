package usecases

import (
	"context"
	"fmt"

	"github.com/dmtr/mail_me_all/backend/errors"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	pb "github.com/dmtr/mail_me_all/backend/rpc"
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
	var user models.User

	newUser := pb.UserToken{UserId: userID, AccessToken: accessToken}
	longToken, err := u.RpcClient.GetAccessToken(context.Background(), &newUser)

	if err != nil {
		return user, NewUseCaseError(err.Error(), errors.CantGetToken)
	}

	if userID != longToken.UserId {
		m := fmt.Sprintf("Users ids do not match %s %s", longToken.UserId, userID)
		return user, NewUseCaseError(m, errors.CantGetToken)
	}

	userInfo, err := u.RpcClient.GetUserInfo(context.Background(), longToken)
	if err != nil {
		return user, NewUseCaseError(err.Error(), errors.CantGetUserInfo)
	}

	var userExists bool
	user, err = u.UserDatastore.GetUserByFbID(ctx, longToken.UserId)
	if err != nil {
		code := errors.GetErrorCode(err)
		if code == errors.NotFound {
			userExists = false
		} else {
			return models.User{}, NewUseCaseError(err.Error(), code)
		}
	} else {
		userExists = true
	}

	user.Name = userInfo.Name
	user.Email = userInfo.Email
	user.FbID = longToken.UserId

	if userExists {
		if user, err = u.UserDatastore.UpdateUser(ctx, user); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
		t := models.Token{
			UserID:  user.ID,
			FbToken: longToken.AccessToken,
		}
		t.ExpiresAt = t.CalculateExpiresAt(longToken.ExpiresIn)

		if _, err := u.UserDatastore.UpdateToken(ctx, t); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
	} else {
		if user, err = u.UserDatastore.InsertUser(ctx, user); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
		log.Debugf("New user %s", user)

		t := models.Token{
			UserID:  user.ID,
			FbToken: longToken.AccessToken,
		}
		t.ExpiresAt = t.CalculateExpiresAt(longToken.ExpiresIn)

		if token, err := u.UserDatastore.InsertToken(ctx, t); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		} else {
			log.Debugf("New token %s", token)
		}
	}
	return user, nil
}

func (u UserUseCase) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	user, err := u.UserDatastore.GetUserByID(ctx, userID)
	if err != nil {
		return user, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}
	return user, err
}
