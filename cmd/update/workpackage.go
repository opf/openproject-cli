package update

import (
	"fmt"
	"strconv"
	
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/workpackages"
	"github.com/spf13/cobra"
)

var actionFlag string

var workPackageCmd = &cobra.Command{
	Use:   "workpackage [id]",
	Short: "Updates the work package",
	Long: `Get a list of unread notifications.
The list can get filtered further.`,
	Run: updateWorkPackage,
}

func updateWorkPackage(_ *cobra.Command, args []string) {
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		printer.ErrorText(fmt.Sprintf("'%s' is an invalid work package id. Must be a number.", args[0]))
	}
	
	workpackages.Update(id)
}
