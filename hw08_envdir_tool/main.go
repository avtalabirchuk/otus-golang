package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Program shud run as: %s <path to directory with env> <command> <some arg>",
			os.Args[0])
	}

	envDir := os.Args[1]
	commandAndArg := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatalf("Unable to get environment variables info: %v", err)
	}
	code := RunCmd(commandAndArg, env)
	os.Exit(code)
}
