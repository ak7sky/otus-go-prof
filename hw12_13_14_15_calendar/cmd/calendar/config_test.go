package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigInit(t *testing.T) {
	configFile = "../../configs/config.yml"
	config, err := NewConfig()
	require.NoError(t, err, "unexpected error")
	require.Equal(t, "INFO", config.Logger.Level)
	require.Equal(t, "memory", config.Storage.Type)
	require.Equal(t, "localhost", config.Server.Host)
	require.Equal(t, 8080, config.Server.HTTPPort)
}

func TestConfigInitFailure(t *testing.T) {
	configFile = "/not-found.yml"
	_, err := NewConfig()
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
