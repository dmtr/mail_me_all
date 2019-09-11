package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dmtr/mail_me_all/backend/config"
	"github.com/dmtr/mail_me_all/backend/db"
	"github.com/dmtr/mail_me_all/backend/models"
	"github.com/dmtr/mail_me_all/backend/usecases"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	retry time.Duration = 4
)

type testFunc func(t *testing.T, router *gin.Engine)

func PerformRequest(r http.Handler, method, path string, body io.Reader, json bool) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	if body != nil && json {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func RunTests(tests map[string]testFunc, t *testing.T) {

	conf := config.GetConfig()
	db_, err := db.ConnectDb(conf.DSN, retry*time.Second)
	if err != nil {
		t.Fatal("Can't connect to database")
	}

	userDatastore := db.NewUserDatastore(db_)
	userUseCase := usecases.NewUserUseCase(userDatastore)
	usecases := models.NewUseCases(userUseCase)
	router := GetRouter(&usecases)

	for name, fn := range tests {
		fmt.Printf("Running test %s", name)
		f := func(t *testing.T) { fn(t, router) }
		t.Run(name, f)
	}
}
