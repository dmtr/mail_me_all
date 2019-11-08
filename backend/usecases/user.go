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
	RpcClient     pb.TwProxyServiceClient
}

// NewUserUseCase implementation
func NewUserUseCase(datastore models.UserDatastore, client pb.TwProxyServiceClient) *UserUseCase {
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
func (u UserUseCase) SignInWithTwitter(ctx context.Context, twitterID, name, email, screenName, accessToken, tokenSecret string) (models.User, error) {
	user := models.User{Name: name}
	var profileUrl string

	req := pb.UserInfoRequest{
		TwitterId:    twitterID,
		AccessToken:  accessToken,
		AccessSecret: tokenSecret,
		ScreenName:   screenName,
	}
	if userInfo, err := u.RpcClient.GetUserInfo(context.Background(), &req); err != nil {
		log.Errorf("Can not get user info: %s", err)
	} else {
		user.Email = userInfo.Email
		profileUrl = userInfo.ProfileImageUrl
	}

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
		twitterUser.ProfileIMGURL = profileUrl

		if _, err = u.UserDatastore.UpdateTwitterUser(ctx, twitterUser); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
	} else {
		if user, err = u.UserDatastore.InsertUser(ctx, user); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}
		log.Debugf("New user %s", user)

		twUser := models.TwitterUser{
			UserID:        user.ID,
			TwitterID:     twitterID,
			AccessToken:   accessToken,
			TokenSecret:   tokenSecret,
			ProfileIMGURL: profileUrl,
		}

		if _, err = u.UserDatastore.InsertTwitterUser(ctx, twUser); err != nil {
			return models.User{}, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
		}

	}
	return user, err

}

func (u UserUseCase) SearchTwitterUsers(ctx context.Context, userID uuid.UUID, query string) ([]models.TwitterUserSearchResult, error) {
	twitterUser, err := u.UserDatastore.GetTwitterUser(ctx, userID)

	if err != nil {
		return nil, err
	}

	req := pb.UserSearchRequest{
		TwitterId:    twitterUser.TwitterID,
		AccessToken:  twitterUser.AccessToken,
		AccessSecret: twitterUser.TokenSecret,
		Query:        query,
	}

	res, err := u.RpcClient.SearchUsers(context.Background(), &req)
	if err != nil {
		log.Errorf("Can not find users: %s", err)
		return nil, err
	}

	users := make([]models.TwitterUserSearchResult, 0, len(res.Users))
	for _, user := range res.Users {
		u := models.TwitterUserSearchResult{
			TwitterID:     user.TwitterId,
			Name:          user.Name,
			ScreenName:    user.ScreenName,
			ProfileIMGURL: user.ProfileImageUrl,
		}
		users = append(users, u)
	}

	return users, err
}

func (u UserUseCase) AddSubscription(ctx context.Context, subscription models.Subscription) (models.Subscription, error) {
	_, err := u.UserDatastore.GetUser(ctx, subscription.UserID)
	if err != nil {
		return subscription, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}

	s, err := u.UserDatastore.InsertSubscription(ctx, subscription)
	if err != nil {
		return subscription, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}

	return s, nil
}

func (u UserUseCase) GetSubscriptions(ctx context.Context, userID uuid.UUID) ([]models.Subscription, error) {
	s, err := u.UserDatastore.GetSubscriptions(ctx, userID)
	if err != nil {
		return s, NewUseCaseError(err.Error(), errors.GetErrorCode(err))
	}

	return s, nil
}
