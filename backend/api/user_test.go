package api

import (
	"bytes"
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
	assert.Equal(t, res.Users[0].TwitterId, r.Users[0].ID)
	assert.Equal(t, res.Users[1].TwitterId, r.Users[1].ID)
}

func testSearchTwitterUsersBadRequest(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	w := performGetRequest(router, "/api/twitter-users", nil)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func testAddSubscriptionUserNotFound(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	e := &db.DbError{Err: sql.ErrNoRows}
	datastoreMock.On("GetUser", mock.Anything, mock.Anything).Return(models.User{}, e)

	req := map[string]interface{}{
		"id": nil, "title": "abc", "email": "test@example.com", "day": "monday",
		"userList": []twitterUser{twitterUser{ID: "123", Name: "test", ScreenName: "test", ProfileIMGURL: "url"}}}
	reqJson, _ := json.Marshal(req)

	w := performPostRequest(router, "/api/subscriptions", bytes.NewBuffer(reqJson))
	assert.Equal(t, http.StatusNotFound, w.Code)

	datastoreMock.AssertNumberOfCalls(t, "GetUser", 1)
	datastoreMock.AssertNumberOfCalls(t, "InsertSubscription", 0)
}

func testUpdateSubscriptionNotFound(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	e := &db.DbError{Err: sql.ErrNoRows}
	datastoreMock.On("UpdateSubscription", mock.Anything, mock.Anything).Return(models.Subscription{}, e)

	email := "test@example.com"
	uid, _ := uuid.Parse(testUserID)
	userEmail := models.UserEmail{
		UserID: uid,
		Email:  email,
		Status: models.EmailStatusConfirmed,
	}
	datastoreMock.On("GetUserEmail", mock.Anything, mock.Anything).Return(userEmail, nil)

	req := map[string]interface{}{
		"id": uuid.New().String(), "title": "abc", "email": email, "day": "monday",
		"userList": []twitterUser{twitterUser{ID: "123", Name: "test", ScreenName: "test", ProfileIMGURL: "url"}}}
	reqJson, _ := json.Marshal(req)

	w := performPutRequest(router, "/api/subscriptions", bytes.NewBuffer(reqJson))
	assert.Equal(t, http.StatusNotFound, w.Code)

	datastoreMock.AssertNumberOfCalls(t, "UpdateSubscription", 1)
}

func testUpdateSubscriptionFailedCantSendEmail(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	title := "abc"
	email := "test@example.com"
	s := models.Subscription{
		ID:    uuid.New(),
		Title: title,
		Email: email,
		Day:   "monday",
	}
	datastoreMock.On("UpdateSubscription", mock.Anything, mock.Anything).Return(s, nil)

	e := &db.DbError{Err: sql.ErrNoRows}
	datastoreMock.On("GetUserEmail", mock.Anything, mock.Anything).Return(models.UserEmail{}, e)

	req := map[string]interface{}{
		"id": uuid.New().String(), "title": title, "email": email, "day": "monday",
		"userList": []twitterUser{twitterUser{ID: "123", Name: "test", ScreenName: "test", ProfileIMGURL: "url"}}}
	reqJson, _ := json.Marshal(req)

	w := performPutRequest(router, "/api/subscriptions", bytes.NewBuffer(reqJson))
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	datastoreMock.AssertNumberOfCalls(t, "GetUserEmail", 1)
	datastoreMock.AssertNumberOfCalls(t, "UpdateSubscription", 0)
}

func testUpdateSubscriptionOk(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	id := uuid.New()
	datastoreMock.On("UpdateSubscription", mock.Anything, mock.Anything).Return(models.Subscription{ID: id}, nil)

	email := "test@example.com"
	uid, _ := uuid.Parse(testUserID)
	userEmail := models.UserEmail{
		UserID: uid,
		Email:  email,
		Status: models.EmailStatusConfirmed,
	}
	datastoreMock.On("GetUserEmail", mock.Anything, mock.Anything).Return(userEmail, nil)

	req := map[string]interface{}{
		"id": id.String(), "title": "abc", "email": email, "day": "monday",
		"userList": []twitterUser{twitterUser{ID: "123", Name: "test", ScreenName: "test", ProfileIMGURL: "url"}}}
	reqJson, _ := json.Marshal(req)

	w := performPutRequest(router, "/api/subscriptions", bytes.NewBuffer(reqJson))
	assert.Equal(t, http.StatusOK, w.Code)

	datastoreMock.AssertNumberOfCalls(t, "UpdateSubscription", 1)
	datastoreMock.AssertNumberOfCalls(t, "UpdateSubscription", 1)
}

func testDeleteSubscriptionNotAuth(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	id, _ := uuid.Parse("1c61dcb2-8bdb-4e3a-8415-d73b1d6133d0")
	s := models.Subscription{
		ID:     uuid.New(),
		UserID: id,
		Title:  "test",
	}
	datastoreMock.On("GetSubscription", mock.Anything, mock.Anything).Return(s, nil)

	req := map[string]interface{}{}
	reqJson, _ := json.Marshal(req)

	w := performDeleteRequest(router, fmt.Sprintf("/api/subscriptions/%s", s.ID.String()), bytes.NewBuffer(reqJson))
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	datastoreMock.AssertNumberOfCalls(t, "GetSubscription", 1)
}

func testDeleteAccountOk(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient) {
	uid, _ := uuid.Parse(testUserID)
	datastoreMock.On("RemoveUser", mock.Anything, uid).Return(nil)

	w := performDeleteRequest(router, "/api/user", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUserEndpoints(t *testing.T) {
	tests := map[string]testFunc{
		"TestGetUserOk":                             testGetUserOk,
		"TestGetUserNotFound":                       testGetUserNotFound,
		"TestSearchTwitterUsersOk":                  testSearchTwitterUsersOk,
		"TestSearchTwitterUsersBadRequest":          testSearchTwitterUsersBadRequest,
		"TestUpdateSubscriptionNotFound":            testUpdateSubscriptionNotFound,
		"TestAddSubscriptionUserNotFound":           testAddSubscriptionUserNotFound,
		"TestDeleteSubscriptionNotAuth":             testDeleteSubscriptionNotAuth,
		"TestDeleteAccountOk":                       testDeleteAccountOk,
		"TestUpdateSubscriptionFailedCantSendEmail": testUpdateSubscriptionFailedCantSendEmail,
	}
	runTests(tests, t)
}
