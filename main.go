package main

import (
	"os"

	"hateb/internal/cli"
	"hateb/internal/hateb"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, hateb.NewClient(nil)))
}
