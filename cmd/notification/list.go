package notification

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/notifications"
)

var listReason string

var validReasons = []string{"", "assigned", "mentioned", "responsible", "watched", "dateAlert"}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists notifications",
	Long: `Get a list of unread notifications.
The list can get filtered further.`,
	Run: listNotifications,
}

func listNotifications(_ *cobra.Command, _ []string) {
	if !common.Contains(validReasons, listReason) {
		printer.ErrorText(fmt.Sprintf("Reason '%s' is invalid.", listReason))
		return
	}

	if all, err := notifications.All(listReason); err == nil {
		printer.Notifications(all)
	} else {
		printer.Error(err)
	}
}
