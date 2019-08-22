package config

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// TestGetConfig func
func TestGetDefaultConfig(t *testing.T) {
	conf := GetConfig()
	assert.Equal(t, appHost, conf.Host, fmt.Sprintf("Host must be %s", appHost))
	assert.Equal(t, appPort, conf.Port, fmt.Sprintf("Port must be %d", appPort))
	assert.Equal(t, log.DebugLevel, conf.Loglevel, "Loglevel must be debug")
}

// TestGetConfig func
func TestGetConfigFromEnv(t *testing.T) {
	os.Setenv(appPrefix+"_HOST", "localhost")
	port := 80
	os.Setenv(appPrefix+"_PORT", strconv.Itoa(port))
	os.Setenv(appPrefix+"_LOGLEVEL", "info")
	conf := GetConfig()
	assert.Equal(t, "localhost", conf.Host, "Host must be localhost")
	assert.Equal(t, port, conf.Port, "Port must be 80")
	assert.Equal(t, log.InfoLevel, conf.Loglevel, "Loglevel must be info")

}
