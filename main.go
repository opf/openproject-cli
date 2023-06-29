package main

import (
	"fmt"
	"os"

	"github.com/opf/openproject-cli/cmd"
	"github.com/opf/openproject-cli/components/configuration"
)

var (
	version = "current"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cliVersion := configuration.BuildCliVersion(version, commit, date)

	if err := cmd.Execute(cliVersion); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
