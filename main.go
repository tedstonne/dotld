package main

import (
	"os"

	"dotld/internal/cli"
)

var version = "dev"

func withImplicitSearch(args []string) []string {
	if len(args) == 0 {
		return args
	}
	first := args[0]
	if first == "search" || first == "--help" || first == "-h" {
		return args
	}

	return append([]string{"search"}, args...)
}

func main() {
	os.Exit(cli.Run(withImplicitSearch(os.Args[1:])))
}
