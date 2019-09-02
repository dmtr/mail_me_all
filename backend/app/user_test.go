package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testCreateUser(t *testing.T, app App) {
	req := map[string]string{"name": "Test Me"}
	req_json, _ := json.Marshal(req)
	w := PerformRequest(app.Router, "POST", "/api/users", bytes.NewBuffer(req_json), true)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	value, exists := response["name"]
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, req["name"], value)
}

func TestUserEndpoinds(t *testing.T) {
	tests := map[string]testFunc{
		"testCreateUser": testCreateUser,
	}
	RunTests(tests, t)
}
