package list

import (
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/status"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Lists status",
	Long:  "Get a list of all status of the instance.",
	Run:   listStatus,
}

func listStatus(_ *cobra.Command, _ []string) {
	if all, err := status.All(); err == nil {
		printer.StatusList(all)
	} else {
		printer.Error(err)
	}
}
