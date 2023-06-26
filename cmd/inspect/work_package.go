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
	Use:   "workpackage [id]",
	Short: "Show details about a work package",
	Long:  "Show detailed information of a work package refereced by it's ID.",
	Run:   inspectWorkPackage,
}

func inspectWorkPackage(_ *cobra.Command, args []string) {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		printer.ErrorText(fmt.Sprintf("'%s' is an invalid work package id. Must be a number.", args[0]))
	}

	workPackage := work_packages.Lookup(id)

	if shouldOpenWorkPackageInBrowser {
		err := launch.Browser(routes.WorkPackageUrl(workPackage))

		if err != nil {
			printer.ErrorText(fmt.Sprintf("Error opening browser: %s", err))
		}
	} else {
		printer.WorkPackage(workPackage)
	}
}
