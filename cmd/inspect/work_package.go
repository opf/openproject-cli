package inspect

import (
	"fmt"
	"github.com/opf/openproject-cli/components/resources/work_packages"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
)

var inspectWorkPackageCmd = &cobra.Command{
	Use:   "workpackage [id]",
	Short: "Show details about a work package",
	Long:  "Show detailed information of a work package refereced by it's ID.",
	Run:   inspectWorkPackage,
}

func inspectWorkPackage(cmd *cobra.Command, args []string) {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		printer.ErrorText(fmt.Sprintf("'%s' is an invalid work package id. Must be a number.", args[0]))
	}

	printer.WorkPackage(work_packages.Lookup(id))
}
