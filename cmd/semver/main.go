package main

import (
	"os"

	"github.com/binary-install/semver/cmd/semver/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
