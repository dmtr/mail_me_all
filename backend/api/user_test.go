package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/dmtr/mail_me_all/backend/mocks"
	"github.com/dmtr/mail_me_all/backend/rpc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testSignUpFBOk(t *testing.T, router *gin.Engine, clientMock *mocks.FbProxyServiceClient) {
	token := rpc.UserToken{UserId: "0011", AccessToken: "2fe", ExpiresIn: 1000}
	clientMock.On("GetAccessToken", mock.Anything, mock.Anything).Return(&token, nil)

	userInfo := rpc.UserInfo{Name: "test", UserId: "0011", Email: "test@example.com"}
	clientMock.On("GetUserInfo", mock.Anything, mock.Anything).Return(&userInfo, nil)

	req := map[string]string{"fbid": token.UserId, "fbtoken": "1abc"}
	reqJson, _ := json.Marshal(req)

	w := PerformRequest(router, "POST", "/api/signin/fb", bytes.NewBuffer(reqJson), true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	value, exists := response["id"]
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, req["fbid"], value)
}

func testSignUpFBFailedBadToken(t *testing.T, router *gin.Engine, clientMock *mocks.FbProxyServiceClient) {
	clientMock.On("GetAccessToken", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("Invalid token"))

	req := map[string]string{"fbid": "000", "fbtoken": "1abc"}
	reqJson, _ := json.Marshal(req)

	w := PerformRequest(router, "POST", "/api/signin/fb", bytes.NewBuffer(reqJson), true)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUserEndpoinds(t *testing.T) {
	tests := map[string]testFunc{
		"testSignUpFB":               testSignUpFBOk,
		"testSignUpFBFailedBadToken": testSignUpFBFailedBadToken,
	}
	RunTests(tests, t)
}
