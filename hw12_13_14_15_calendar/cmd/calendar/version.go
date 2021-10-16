package main

import (
	"flag"
	"fmt"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func printVersion() {
	fmt.Printf("Calendar %s release (%s) built on %s\n", release, gitHash, buildDate)
}

func isVersionCommand() bool {
	args := flag.Args()
	if len(args) != 0 && args[0] == "version" {
		return true
	}
	return false
}
