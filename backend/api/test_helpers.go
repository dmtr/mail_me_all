package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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
	router := GetRouter()
	for name, fn := range tests {
		fmt.Printf("Running test %s", name)
		f := func(t *testing.T) { fn(t, router) }
		t.Run(name, f)
	}
}
