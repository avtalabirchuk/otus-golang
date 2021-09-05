package main

import (
	"os"
)

func main() {
	// programName := os.Args[0]

	// if len(os.Args) < 3 {
	// 	log.Fatalf("Program shud run as: %s <path to directory with env> <command> <some arg>", programName)
	// }

	envDir := os.Args[1]
	commandAndArg := os.Args[2:]

	env, _ := ReadDir(envDir)
	// if err != nil {
	// 	log.Fatalf("Unable to get environment variables info: %v", err)
	// }
	code := RunCmd(commandAndArg, env)
	os.Exit(code)
}
