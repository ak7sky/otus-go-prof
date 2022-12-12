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
	if err != nil {
		require.Fail(t, "error creating file for cmd out")
	}
	defer func() {
		if err := os.Remove("testdata/cmdOutFile"); err != nil {
			require.Fail(t, "error removing cmd out file after test")
		}
	}()

	os.Stdout = cmdOutFile
	defer func() { os.Stdout = originalStdOut }()

	setOsEnv(t)

	cmd := []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}

	returnCode := RunCmd(cmd, expectedEnv)
	cmdOutBytes, err := os.ReadFile("testdata/cmdOutFile")
	if err != nil {
		require.Fail(t, "error during read teat data")
	}

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
	if err := os.Setenv("HELLO", "SHOULD_REPLACE"); err != nil {
		require.Fail(t, "error preparing test environment")
	}
	if err := os.Setenv("FOO", "SHOULD_REPLACE"); err != nil {
		require.Fail(t, "error preparing test environment")
	}
	if err := os.Setenv("UNSET", "SHOULD_REMOVE"); err != nil {
		require.Fail(t, "error preparing test environment")
	}
	if err := os.Setenv("ADDED", "from original env"); err != nil {
		require.Fail(t, "error preparing test environment")
	}
	if err := os.Setenv("EMPTY", "SHOULD_BE_EMPTY"); err != nil {
		require.Fail(t, "error preparing test environment")
	}
}
