package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/dmtr/mail_me_all/backend/db"
	"github.com/dmtr/mail_me_all/backend/mocks"
	"github.com/dmtr/mail_me_all/backend/models"
	pb "github.com/dmtr/mail_me_all/backend/rpc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testGetUserOk(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	uid, _ := uuid.Parse(testUserID)
	user := models.User{ID: uid, Name: "test"}
	datastoreMock.On("GetUser", mock.Anything, uid).Return(user, nil)

	w := performGetRequest(router, "/api/user", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var res appUser
	err := json.Unmarshal([]byte(w.Body.String()), &res)
	assert.NoError(t, err)
	assert.Equal(t, testUserID, res.ID)
	assert.Equal(t, true, res.SignedIn)
}

func testGetUserNotFound(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	uid, _ := uuid.Parse(testUserID)
	e := &db.DbError{Err: sql.ErrNoRows}
	datastoreMock.On("GetUser", mock.Anything, uid).Return(models.User{}, e)

	w := performGetRequest(router, "/api/user", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func testSearchTwitterUsersOk(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	uid, _ := uuid.Parse(testUserID)
	user := models.TwitterUser{UserID: uid, TwitterID: "111", AccessToken: "token", TokenSecret: "secret"}
	datastoreMock.On("GetTwitterUser", mock.Anything, uid).Return(user, nil)

	req := pb.UserSearchRequest{
		TwitterId:    user.TwitterID,
		AccessToken:  user.AccessToken,
		AccessSecret: user.TokenSecret,
		Query:        "test",
	}

	res := pb.UserSearchResult{
		Users: []*pb.UserInfo{&pb.UserInfo{TwitterId: "222", ScreenName: "foo", ProfileImageUrl: "some_url"},
			&pb.UserInfo{TwitterId: "333", ScreenName: "bar", ProfileImageUrl: "other_url"},
		},
	}
	clientMock.On("SearchUsers", mock.Anything, &req).Return(&res, nil)

	w := performGetRequest(router, "/api/twitter-users?q=test", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var r struct {
		Users []twitterUser
	}
	err := json.Unmarshal([]byte(w.Body.String()), &r)
	assert.NoError(t, err)
	fmt.Printf("res %v", r)
	assert.Equal(t, 2, len(r.Users))

}

func testSearchTwitterUsersBadRequest(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	w := performGetRequest(router, "/api/twitter-users", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestUserEndpoints(t *testing.T) {
	tests := map[string]testFunc{
		"TestGetUserOk":                    testGetUserOk,
		"TestGetUserNotFound":              testGetUserNotFound,
		"TestSearchTwitterUsersOk":         testSearchTwitterUsersOk,
		"TestSearchTwitterUsersBadRequest": testSearchTwitterUsersBadRequest,
	}
	runTests(tests, t)
}
