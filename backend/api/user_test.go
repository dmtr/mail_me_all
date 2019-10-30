package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/dmtr/mail_me_all/backend/mocks"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testGetUserOk(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore) {
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

func TestUserEndpoints(t *testing.T) {
	tests := map[string]testFunc{
		"TestGetUserOk": testGetUserOk,
	}
	runTests(tests, t)
}
