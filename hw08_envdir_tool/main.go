package main

import (
	"log"
	"os"
)

func main() {
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(RunCmd(os.Args[2:], env))
}
