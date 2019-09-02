package app

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testFunc func(t *testing.T, app App)

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
	app := GetApp()
	defer app.Close()

	for name, fn := range tests {
		fmt.Printf("Running test %s", name)
		f := func(t *testing.T) { fn(t, app) }
		t.Run(name, f)
	}
}
