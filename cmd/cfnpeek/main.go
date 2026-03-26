package main

import (
	"fmt"
	"os"

	"github.com/awesome-foundation/cfnpeek/internal/cli"
)

// Set by goreleaser ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd := cli.NewRootCmd(fmt.Sprintf("%s (%s, %s)", version, commit, date))
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
