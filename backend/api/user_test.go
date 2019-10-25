package api

import (
	"testing"
)

//func testSignUpWithTwitterOk(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore) {
//	token := rpc.UserToken{UserId: "0011", AccessToken: "2fe", ExpiresIn: 1000}
//	clientMock.On("GetAccessToken", mock.Anything, mock.Anything).Return(&token, nil)
//
//	userInfo := rpc.UserInfo{Name: "test", UserId: "0011", Email: "test@example.com"}
//	clientMock.On("GetUserInfo", mock.Anything, mock.Anything).Return(&userInfo, nil)
//
//	e := &db.DbError{Err: sql.ErrNoRows}
//	datastoreMock.On("GetUserByFbID", mock.Anything, mock.Anything).Return(models.User{}, e)
//
//	uid := uuid.New()
//	user := models.User{ID: uid, Name: userInfo.Name, Email: userInfo.Email, FbID: userInfo.UserId}
//	datastoreMock.On("InsertUser", mock.Anything, mock.Anything).Return(user, nil)
//
//	tokenDb := models.Token{
//		UserID: uid, FbToken: token.AccessToken, ExpiresAt: models.CalculateExpiresAt(token.ExpiresIn)}
//	datastoreMock.On("InsertToken", mock.Anything, mock.Anything).Return(tokenDb, nil)
//
//	req := map[string]string{"fbid": token.UserId, "fbtoken": "1abc"}
//	reqJson, _ := json.Marshal(req)
//
//	w := PerformPostRequest(router, "/api/signin/fb", bytes.NewBuffer(reqJson))
//	assert.Equal(t, http.StatusOK, w.Code)
//	h := w.Header()
//	c := h.Get("set-cookie")
//	contains := strings.Contains(c, "session")
//	assert.True(t, contains)
//
//	clientMock.AssertNumberOfCalls(t, "GetAccessToken", 1)
//	clientMock.AssertNumberOfCalls(t, "GetUserInfo", 1)
//
//	datastoreMock.AssertNumberOfCalls(t, "InsertUser", 1)
//	datastoreMock.AssertNumberOfCalls(t, "InsertToken", 1)
//
//	var response appUser
//	err := json.Unmarshal([]byte(w.Body.String()), &response)
//	assert.Nil(t, err)
//	assert.Equal(t, true, response.SignedIn)
//	assert.NotEmpty(t, response.ID)
//
//	datastoreMock.On("GetUserByID", mock.Anything, mock.Anything).Return(user, nil)
//
//	w = PerformGetRequest(router, "/api/user", ParseCookie(c))
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	var res appUser
//	err = json.Unmarshal([]byte(w.Body.String()), &res)
//	assert.Equal(t, response.ID, res.ID)
//	assert.Equal(t, true, res.SignedIn)
//}

//func testSignInFBOk(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore) {
//	token := rpc.UserToken{UserId: "0011", AccessToken: "2fe", ExpiresIn: 1000}
//	clientMock.On("GetAccessToken", mock.Anything, mock.Anything).Return(&token, nil)
//
//	userInfo := rpc.UserInfo{Name: "test", UserId: "0011", Email: "test@example.com"}
//	clientMock.On("GetUserInfo", mock.Anything, mock.Anything).Return(&userInfo, nil)
//
//	uid := uuid.New()
//	user := models.User{ID: uid, Name: userInfo.Name, Email: "old@email.me", FbID: userInfo.UserId}
//	datastoreMock.On("GetUserByFbID", mock.Anything, token.UserId).Return(user, nil)
//
//	updatedUser := models.User{ID: uid, Name: user.Name, Email: userInfo.Email, FbID: userInfo.UserId}
//	datastoreMock.On("UpdateUser", mock.Anything, updatedUser).Return(updatedUser, nil)
//
//	updatedToken := models.Token{
//		UserID: uid, FbToken: token.AccessToken, ExpiresAt: models.CalculateExpiresAt(token.ExpiresIn)}
//	datastoreMock.On("UpdateToken", mock.Anything, mock.Anything).Return(updatedToken, nil)
//
//	req := map[string]string{"fbid": token.UserId, "fbtoken": "1abc"}
//	reqJson, _ := json.Marshal(req)
//
//	w := PerformPostRequest(router, "/api/signin/fb", bytes.NewBuffer(reqJson))
//	assert.Equal(t, http.StatusOK, w.Code)
//	h := w.Header()
//	c := h.Get("set-cookie")
//	contains := strings.Contains(c, "session")
//	assert.True(t, contains)
//
//	clientMock.AssertNumberOfCalls(t, "GetAccessToken", 1)
//	clientMock.AssertNumberOfCalls(t, "GetUserInfo", 1)
//
//	datastoreMock.AssertNumberOfCalls(t, "GetUserByFbID", 1)
//	datastoreMock.AssertNumberOfCalls(t, "UpdateUser", 1)
//	datastoreMock.AssertNumberOfCalls(t, "GetToken", 0)
//	datastoreMock.AssertNumberOfCalls(t, "UpdateToken", 1)
//
//	var response appUser
//	err := json.Unmarshal([]byte(w.Body.String()), &response)
//	assert.Nil(t, err)
//	assert.NotEmpty(t, response.ID)
//	assert.Equal(t, true, response.SignedIn)
//}

func TestUserEndpoints(t *testing.T) {
	tests := map[string]testFunc{}
	runTests(tests, t)
}
