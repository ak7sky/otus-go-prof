package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedEnv = Environment{
	"BAR":   EnvValue{Value: "bar"},
	"EMPTY": EnvValue{},
	"FOO":   EnvValue{Value: "   foo\nwith new line"},
	"HELLO": EnvValue{Value: `"hello"`},
	"UNSET": EnvValue{NeedRemove: true},
}

func TestReadDir(t *testing.T) {
	actualEnv, err := ReadDir("testdata/env")
	require.NoError(t, err)
	require.Equal(t, expectedEnv, actualEnv)
}

func TestReadDirFailure(t *testing.T) {
	_, err := os.Create("testdata/env/FOO=BAR")
	require.NoError(t, err, "error during creating test file")

	defer func() {
		err := os.Remove("testdata/env/FOO=BAR")
		require.NoError(t, err, "error during removing file after test")
	}()

	actualEnv, err := ReadDir("testdata/env")
	require.Nil(t, actualEnv)
	require.EqualError(t, err, `env file name contains "="`)
}
