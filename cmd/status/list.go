package status

import (
	"github.com/spf13/cobra"

	openerrors "github.com/opf/openproject-cli/components/errors"
	"github.com/opf/openproject-cli/components/printer"
	wpstatus "github.com/opf/openproject-cli/components/resources/status"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists statuses",
	Long:  "Get a list of all statuses of the instance.",
	RunE:  listStatuses,
}

func listStatuses(_ *cobra.Command, _ []string) error {
	if all, err := wpstatus.All(); err == nil {
		printer.StatusList(all)
	} else {
		printer.Error(err)
		return openerrors.ErrHandled
	}
	return nil
}
