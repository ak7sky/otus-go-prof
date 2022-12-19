package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	exCmd := cmd[0]
	cmdWithArgs := exec.Command(exCmd, cmd[1:]...)
	cmdWithArgs.Stdin = os.Stdin
	cmdWithArgs.Stdout = os.Stdout
	cmdWithArgs.Stderr = os.Stderr

	for name, ev := range env {
		if ev.NeedRemove {
			if err := os.Unsetenv(name); err != nil {
				log.Println(err)
				return 1
			}
			continue
		}
		if err := os.Setenv(name, ev.Value); err != nil {
			log.Println(err)
			return 1
		}
	}
	cmdWithArgs.Env = os.Environ()

	if err := cmdWithArgs.Run(); err != nil {
		log.Println(err.Error())
	}

	return cmdWithArgs.ProcessState.ExitCode()
}
