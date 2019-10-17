package usecases

import (
	"context"
	"fmt"

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
		return user, NewUseCaseError(err.Error(), cantGetToken)
	}

	if userID != longToken.UserId {
		m := fmt.Sprintf("Users ids do not match %s %s", longToken.UserId, userID)
		return user, NewUseCaseError(m, cantGetToken)
	}

	userInfo, err := u.RpcClient.GetUserInfo(context.Background(), longToken)
	if err != nil {
		return user, NewUseCaseError(err.Error(), cantGetUserInfo)
	}

	var userExists bool
	user, err = u.UserDatastore.GetUserByFbID(ctx, longToken.UserId)
	if err != nil {
		code := getErrorCode(err)
		if code == notFound {
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
			return models.User{}, NewUseCaseError(err.Error(), getErrorCode(err))
		}
		var t models.Token
		t, err = u.UserDatastore.GetToken(ctx, user.ID)
		if err != nil {
			return models.User{}, NewUseCaseError(err.Error(), getErrorCode(err))
		}

		if _, err := u.UserDatastore.UpdateToken(ctx, t); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), getErrorCode(err))
		}
	} else {
		if user, err = u.UserDatastore.InsertUser(ctx, user); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), getErrorCode(err))
		}
		log.Debugf("New user %s", user)

		t := models.Token{
			UserID:  user.ID,
			FbToken: longToken.AccessToken,
		}
		t.ExpiresAt = t.CalculateExpiresAt(longToken.ExpiresIn)

		if token, err := u.UserDatastore.InsertToken(ctx, t); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), getErrorCode(err))
		} else {
			log.Debugf("New token %s", token)
		}
	}
	return user, nil
}

func (u UserUseCase) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	user, err := u.UserDatastore.GetUserByID(ctx, userID)
	if err != nil {
		return user, NewUseCaseError(err.Error(), getErrorCode(err))
	}
	return user, err
}
