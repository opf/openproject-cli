package status

import (
	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	wpstatus "github.com/opf/openproject-cli/components/resources/status"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists statuses",
	Long:  "Get a list of all statuses of the instance.",
	Run:   listStatuses,
}

func listStatuses(_ *cobra.Command, _ []string) {
	if all, err := wpstatus.All(); err == nil {
		printer.StatusList(all)
	} else {
		printer.Error(err)
	}
}
