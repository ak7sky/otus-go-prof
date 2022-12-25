package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigInit(t *testing.T) {
	dsn := "postgres://test_user:test_pswd@localhost/calendar"
	err := os.Setenv("DSN", dsn)
	require.NoError(t, err)

	config, err := NewConfig("../../configs/config.yml")
	require.NoError(t, err, "unexpected error")

	require.Equal(t, "INFO", config.Logger.Level)
	require.True(t, config.Logger.IsJSONEnabled)
	require.Equal(t, "memory", config.Storage.Type)
	require.Equal(t, "localhost", config.Server.Host)
	require.Equal(t, 8080, config.Server.HTTPPort)
	require.Equal(t, dsn, config.Storage.DSN)
}

func TestConfigInitFailure(t *testing.T) {
	_, err := NewConfig("./not-found.yml")
	require.Error(t, err, "config init failure expected")
	require.ErrorContains(t, err, "failed to init config")
	fmt.Println(err)
}

func TestConfigValidation(t *testing.T) {
	config := Config{}
	config.Logger.Level = "invalid"
	config.Storage.Type = "invalid"
	config.Server.HTTPPort = -1
	err := validate(config)
	require.Error(t, err)
	require.ErrorContains(t, err, constraints[loggerLevel])
	require.ErrorContains(t, err, constraints[storageType])
	require.ErrorContains(t, err, constraints[serverHTTPPort])
}
