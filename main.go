package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/opf/openproject-cli/cmd"
	"github.com/opf/openproject-cli/components/configuration"
	openerrors "github.com/opf/openproject-cli/components/errors"
)

var (
	version = "current"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cliVersion := configuration.BuildCliVersion(version, commit, date)

	if err := cmd.Execute(cliVersion); err != nil {
		if !errors.Is(err, openerrors.ErrHandled) {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
