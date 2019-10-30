package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/mocks"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/dmtr/mail_me_all/backend/usecases"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const testUserID string = "15b24dd0-1f38-4e0a-8d6f-8df509051279"

type testFunc func(t *testing.T, router *gin.Engine, datastoreMock *mocks.UserDatastore, clientMock *mocks.TwProxyServiceClient)

func performRequest(r http.Handler, method, path string, body io.Reader, json bool, cookie *http.Cookie) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	if body != nil && json {
		req.Header.Set("Content-Type", "application/json")
	}

	if cookie != nil {
		req.AddCookie(cookie)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performPostRequest(r http.Handler, path string, body io.Reader) *httptest.ResponseRecorder {
	return performRequest(r, "POST", path, body, true, nil)
}

func performGetRequest(r http.Handler, path string, cookie *http.Cookie) *httptest.ResponseRecorder {
	return performRequest(r, "GET", path, nil, false, cookie)
}

func parseCookie(cookie string) *http.Cookie {
	c := http.Cookie{}
	s := strings.Split(cookie, ";")
	for _, v := range s {
		keys := strings.Split(v, "=")
		if len(keys) == 2 && keys[0] == "session" {
			c.Name = keys[0]
			c.Value = keys[1]
			break
		}
	}
	return &c
}

func runTests(tests map[string]testFunc, t *testing.T) {

	conf := config.GetConfig()
	conf.Testing = true

	datastoreMock := new(mocks.UserDatastore)
	clientMock := new(mocks.TwProxyServiceClient)
	userUseCase := usecases.NewUserUseCase(datastoreMock, clientMock)
	usecases := models.NewUseCases(userUseCase)
	router := GetRouter(&conf, nil, usecases)

	for name, fn := range tests {
		fmt.Printf("Running test %s", name)
		f := func(t *testing.T) {
			datastoreMock := new(mocks.UserDatastore)
			userUseCase.UserDatastore = datastoreMock

			clientMock = new(mocks.TwProxyServiceClient)
			userUseCase.RpcClient = clientMock

			fn(t, router, datastoreMock, clientMock)
		}
		t.Run(name, f)
	}
}
