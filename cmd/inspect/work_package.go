package inspect

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/launch"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"github.com/opf/openproject-cli/components/routes"
)

var shouldOpenWorkPackageInBrowser bool

var inspectWorkPackageCmd = &cobra.Command{
	Use:     "workpackage [id]",
	Aliases: []string{"wp"},
	Short:   "Show details about a work package",
	Long:    "Show detailed information of a work package referenced by it's ID.",
	Run:     inspectWorkPackage,
}

func inspectWorkPackage(_ *cobra.Command, args []string) {
	if len(args) != 1 {
		printer.ErrorText(fmt.Sprintf("Expected 1 argument [id], but got %d", len(args)))
		return
	}

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		printer.ErrorText(fmt.Sprintf("'%s' is an invalid work package id. Must be a number.", args[0]))
		return
	}

	workPackage, err := work_packages.Lookup(id)
	if err != nil {
		printer.Error(err)
		return
	}

	if shouldOpenWorkPackageInBrowser {
		err = launch.Browser(routes.WorkPackageUrl(workPackage))
		if err != nil {
			printer.ErrorText(fmt.Sprintf("Error opening browser: %+v", err))
		}
	} else {
		printer.WorkPackage(workPackage)
	}
}
