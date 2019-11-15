package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/db"
	"github.com/dmtr/mail_me_all/backend/mocks"
	"github.com/dmtr/mail_me_all/backend/models"
	pb "github.com/dmtr/mail_me_all/backend/rpc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testFunc func(t *testing.T, usecases *models.UseCases, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient)

func runTests(tests map[string]testFunc, t *testing.T) {

	conf := config.GetConfig()
	conf.Testing = true

	datastoreMock := new(mocks.UserDatastore)
	clientMock := new(mocks.TwProxyServiceClient)
	userUseCase := NewUserUseCase(datastoreMock, clientMock)
	systemUseCase := NewSystemUseCase(datastoreMock, clientMock)
	usecases := models.NewUseCases(userUseCase, systemUseCase)

	for name, fn := range tests {
		fmt.Printf("Running test %s", name)
		f := func(t *testing.T) {
			datastoreMock := new(mocks.UserDatastore)
			userUseCase.UserDatastore = datastoreMock

			clientMock := new(mocks.TwProxyServiceClient)
			userUseCase.RpcClient = clientMock

			fn(t, usecases, datastoreMock, clientMock)
		}
		t.Run(name, f)
	}
}

func testSignUpWithTwitterOk(t *testing.T, usecases *models.UseCases, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	twitterUserID := "123"
	screenName := "test"
	accessToken := "access token"
	tokenSecret := "token secret"
	name := "Test"
	email := "test@example.com"

	req := pb.UserInfoRequest{
		TwitterId:    twitterUserID,
		AccessToken:  accessToken,
		AccessSecret: tokenSecret,
		ScreenName:   screenName,
	}

	res := pb.UserInfo{
		TwitterId:       twitterUserID,
		Name:            name,
		Email:           email,
		ScreenName:      screenName,
		ProfileImageUrl: "some_url",
	}
	clientMock.On("GetUserInfo", mock.Anything, &req).Return(&res, nil)

	e := &db.DbError{Err: sql.ErrNoRows}
	datastoreMock.On("GetTwitterUserByID", mock.Anything, twitterUserID).Return(models.TwitterUser{}, e)

	uid := uuid.New()
	user := models.User{ID: uid, Name: name, Email: email}
	datastoreMock.On("InsertUser", mock.Anything, mock.Anything).Return(user, nil)

	twitterUser := models.TwitterUser{
		UserID:        uid,
		TwitterID:     twitterUserID,
		AccessToken:   accessToken,
		TokenSecret:   tokenSecret,
		ProfileIMGURL: res.ProfileImageUrl,
	}
	datastoreMock.On("InsertTwitterUser", mock.Anything, twitterUser).Return(twitterUser, nil)

	u, err := usecases.SignInWithTwitter(context.Background(), twitterUserID, name, email, screenName, accessToken, tokenSecret)
	assert.NoError(t, err)
	assert.Equal(t, user, u)

	datastoreMock.AssertNumberOfCalls(t, "GetTwitterUserByID", 1)
	datastoreMock.AssertNumberOfCalls(t, "InsertUser", 1)
	datastoreMock.AssertNumberOfCalls(t, "InsertTwitterUser", 1)

	clientMock.AssertNumberOfCalls(t, "GetUserInfo", 1)
}

func testSignInWithTwitterOk(t *testing.T, usecases *models.UseCases, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	twitterUserID := "123"
	tokenSecret := "token secret"
	accessToken := "access token"
	uid := uuid.New()
	name := "Test"
	email := "test@example.com"
	screenName := "test"

	req := pb.UserInfoRequest{
		TwitterId:    twitterUserID,
		AccessToken:  accessToken,
		AccessSecret: tokenSecret,
		ScreenName:   screenName,
	}

	res := pb.UserInfo{
		TwitterId:       twitterUserID,
		Name:            name,
		Email:           email,
		ScreenName:      screenName,
		ProfileImageUrl: "some_url",
	}
	clientMock.On("GetUserInfo", mock.Anything, &req).Return(&res, nil)

	twitterUser := models.TwitterUser{
		UserID:        uid,
		TwitterID:     twitterUserID,
		AccessToken:   "old token",
		TokenSecret:   "old secret",
		ProfileIMGURL: res.ProfileImageUrl,
	}

	datastoreMock.On("GetTwitterUserByID", mock.Anything, twitterUserID).Return(twitterUser, nil)

	user := models.User{ID: uid, Name: name, Email: email}
	datastoreMock.On("UpdateUser", mock.Anything, user).Return(user, nil)

	updatedTwitterUser := models.TwitterUser{
		UserID:        twitterUser.UserID,
		TwitterID:     twitterUser.TwitterID,
		AccessToken:   accessToken,
		TokenSecret:   tokenSecret,
		ProfileIMGURL: res.ProfileImageUrl,
	}
	datastoreMock.On("UpdateTwitterUser", mock.Anything, updatedTwitterUser).Return(updatedTwitterUser, nil)

	u, err := usecases.SignInWithTwitter(context.Background(), twitterUserID, name, email, screenName, accessToken, tokenSecret)
	assert.NoError(t, err)
	assert.Equal(t, user, u)

	datastoreMock.AssertNumberOfCalls(t, "GetTwitterUserByID", 1)
	datastoreMock.AssertNumberOfCalls(t, "UpdateTwitterUser", 1)
	datastoreMock.AssertNumberOfCalls(t, "UpdateUser", 1)

	clientMock.AssertNumberOfCalls(t, "GetUserInfo", 1)
}

func TestUseCases(t *testing.T) {
	tests := map[string]testFunc{
		"TestSignUpWithTwitterOk": testSignUpWithTwitterOk,
		"TestSignInWithTwitterOk": testSignInWithTwitterOk,
	}
	runTests(tests, t)
}
