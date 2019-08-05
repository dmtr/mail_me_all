package config

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetConfig func
func TestGetDefaultConfig(t *testing.T) {
	conf := GetConfig()
	assert.Equal(t, appHost, conf.Host, fmt.Sprintf("Host must be %s", appHost))
	assert.Equal(t, appPort, conf.Port, fmt.Sprintf("Port must be %d", appPort))

}

// TestGetConfig func
func TestGetConfigFromEnv(t *testing.T) {
	os.Setenv(appPrefix+"_HOST", "localhost")
	port := 80
	os.Setenv(appPrefix+"_PORT", strconv.Itoa(port))
	conf := GetConfig()
	assert.Equal(t, "localhost", conf.Host, "Host must be localhost")
	assert.Equal(t, port, conf.Port, "Port must be 80")

}
