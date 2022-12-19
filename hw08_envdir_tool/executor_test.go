package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var expSuccessfulOut = `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2
`

func TestRunCmd(t *testing.T) {
	originalStdOut := os.Stdout
	cmdOutFile, err := os.Create("testdata/cmdOutFile")
	require.NoError(t, err, "error creating file for cmd out")
	defer func() {
		err := os.Remove("testdata/cmdOutFile")
		require.NoError(t, err, "error removing cmd out file after test")
	}()

	os.Stdout = cmdOutFile
	defer func() { os.Stdout = originalStdOut }()

	setOsEnv(t)

	cmd := []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}

	returnCode := RunCmd(cmd, expectedEnv)
	cmdOutBytes, err := os.ReadFile("testdata/cmdOutFile")
	require.NoError(t, err, "error during read cmd out from file")

	require.Equal(t, 0, returnCode)
	require.Equal(t, expSuccessfulOut, string(cmdOutBytes))
}

func TestRunCmdFailure(t *testing.T) {
	testCases := []struct {
		name      string
		cmd       []string
		expStatus int
	}{
		{name: "exit status 1", cmd: []string{"/bin/cat", "testdata/unknown.sh"}, expStatus: 1},
		{name: "exit status -1", cmd: []string{"/bin/unknown", "arg1=1"}, expStatus: -1},
		{name: "exit status 127", cmd: []string{"/bin/bash", "testdata/unknown.sh"}, expStatus: 127},
		{name: "exit status 126", cmd: []string{"/bin/bash", "/bin/ls"}, expStatus: 126},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			returnCode := RunCmd(tc.cmd, expectedEnv)
			require.Equal(t, tc.expStatus, returnCode)
		})
	}
}

func setOsEnv(t *testing.T) {
	t.Helper()

	set := func(envName, envVal string) {
		err := os.Setenv(envName, envVal)
		require.NoError(t, err, "error preparing test environment")
	}

	set("HELLO", "SHOULD_REPLACE")
	set("FOO", "SHOULD_REPLACE")
	set("UNSET", "SHOULD_REMOVE")
	set("ADDED", "from original env")
	set("EMPTY", "SHOULD_BE_EMPTY")
}
