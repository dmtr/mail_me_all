package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/mocks"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/dmtr/mail_me_all/backend/usecases"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	retry time.Duration = 4
)

type testFunc func(t *testing.T, router *gin.Engine, clientMock *mocks.FbProxyServiceClient, datastoreMock *mocks.UserDatastore)

func PerformRequest(r http.Handler, method, path string, body io.Reader, json bool, cookie *http.Cookie) *httptest.ResponseRecorder {
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

func PerformPostRequest(r http.Handler, path string, body io.Reader) *httptest.ResponseRecorder {
	return PerformRequest(r, "POST", path, body, true, nil)
}

func PerformGetRequest(r http.Handler, path string, cookie *http.Cookie) *httptest.ResponseRecorder {
	return PerformRequest(r, "GET", path, nil, false, cookie)
}

func ParseCookie(cookie string) *http.Cookie {
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

func RunTests(tests map[string]testFunc, t *testing.T) {

	conf := config.GetConfig()
	conf.Testing = true

	datastoreMock := new(mocks.UserDatastore)
	clientMock := new(mocks.FbProxyServiceClient)
	userUseCase := usecases.NewUserUseCase(datastoreMock, clientMock)
	usecases := models.NewUseCases(userUseCase)
	router := GetRouter(&conf, nil, usecases)

	for name, fn := range tests {
		fmt.Printf("Running test %s", name)
		f := func(t *testing.T) {
			datastoreMock := new(mocks.UserDatastore)
			userUseCase.UserDatastore = datastoreMock

			clientMock = new(mocks.FbProxyServiceClient)
			userUseCase.RpcClient = clientMock

			fn(t, router, clientMock, datastoreMock)
		}
		t.Run(name, f)
	}
}
