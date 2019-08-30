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
	assert.Equal(t, "", conf.DSN, "DSN must be empty string")
}

// TestGetConfig func
func TestGetConfigFromEnv(t *testing.T) {
	os.Setenv(appPrefix+"_HOST", "localhost")
	port := 80
	os.Setenv(appPrefix+"_PORT", strconv.Itoa(port))
	os.Setenv(appPrefix+"_LOGLEVEL", "info")
	dsn := "postgres://postgres@localhost"
	os.Setenv(appPrefix+"_DSN", dsn)

	conf := GetConfig()
	assert.Equal(t, "localhost", conf.Host, "Host must be localhost")
	assert.Equal(t, port, conf.Port, "Port must be 80")
	assert.Equal(t, log.InfoLevel, conf.Loglevel, "Loglevel must be info")
	assert.Equal(t, dsn, conf.DSN, "DSN must be postgres://postgres@localhost")

}
