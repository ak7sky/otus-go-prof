package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, de := range dirEntries {
		if de.IsDir() {
			continue
		}

		if strings.Contains(de.Name(), "=") {
			return nil, errors.New(`env file name contains "="`)
		}

		valueLine, err := readLine(fmt.Sprintf("%s/%s", dir, de.Name()))
		if err != nil {
			return nil, err
		}

		if len(valueLine) == 0 {
			env[de.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		value := strings.TrimRight(valueLine, " \t")
		value = strings.ReplaceAll(value, "\x00", "\n")

		env[de.Name()] = EnvValue{Value: value}
	}

	return env, nil
}

func readLine(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer closeWithCheck(f)

	reader := bufio.NewReader(f)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	return strings.TrimRight(line, "\n"), nil
}

func closeWithCheck(f *os.File) {
	if clErr := f.Close(); clErr != nil {
		log.Fatalf("problem during closing %s; details: %s\n", f.Name(), clErr.Error())
	}
}
