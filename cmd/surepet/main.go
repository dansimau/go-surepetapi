package main

import (
	"os"

	"github.com/dansimau/go-surepetapi/pkg/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
