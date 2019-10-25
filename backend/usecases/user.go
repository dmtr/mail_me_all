package usecases

import (
	"context"

	"github.com/dmtr/mail_me_all/backend/errors"
	"github.com/dmtr/mail_me_all/backend/models"
	pb "github.com/dmtr/mail_me_all/backend/rpc"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// UserUseCase implementation
type UserUseCase struct {
	UserDatastore models.UserDatastore
	RpcClient     pb.FbProxyServiceClient
}

// NewUserUseCase implementation
func NewUserUseCase(datastore models.UserDatastore, client pb.FbProxyServiceClient) *UserUseCase {
	return &UserUseCase{UserDatastore: datastore, RpcClient: client}
}

// GetUserByID implementation
func (u UserUseCase) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	user, err := u.UserDatastore.GetUser(ctx, userID)
	if err != nil {
		return user, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}
	return user, err
}

// SignInWithTwitter implementation
func (u UserUseCase) SignInWithTwitter(ctx context.Context, twitterID, name, email, accessToken, tokenSecret string) (models.User, error) {
	var user models.User
	user.Name = name
	user.Email = email

	var userExists bool

	twitterUser, err := u.UserDatastore.GetTwitterUserByID(ctx, twitterID)
	if err != nil {
		code := errors.GetErrorCode(err)
		if code == errors.NotFound {
			userExists = false
		} else {
			return user, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
	} else {
		userExists = true
	}

	if userExists {
		user.ID = twitterUser.UserID

		if user, err = u.UserDatastore.UpdateUser(ctx, user); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}

		twitterUser.AccessToken = accessToken
		twitterUser.TokenSecret = tokenSecret

		if _, err = u.UserDatastore.UpdateTwitterUser(ctx, twitterUser); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
	} else {
		if user, err = u.UserDatastore.InsertUser(ctx, user); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
		log.Debugf("New user %s", user)

		twUser := models.TwitterUser{
			UserID:      user.ID,
			TwitterID:   twitterID,
			AccessToken: accessToken,
			TokenSecret: tokenSecret,
		}

		if _, err = u.UserDatastore.InsertTwitterUser(ctx, twUser); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}

	}
	return user, err

}
