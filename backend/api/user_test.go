package api

import (
	"net/http"
	"testing"

	"github.com/dmtr/mail_me_all/backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func testGetUserNoCookie(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore) {
	w := performGetRequest(router, "/api/user", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserEndpoints(t *testing.T) {
	tests := map[string]testFunc{}
	runTests(tests, t)
}
