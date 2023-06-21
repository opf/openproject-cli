package list

import (
	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/notifications"
	"github.com/opf/openproject-cli/models/types"
)

var NotificationReason string
var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Lists notifications",
	Long: `Get a list of unread notifications.
The list can get filtered further.`,
	Run: listNotifications,
}

func listNotifications(_ *cobra.Command, _ []string) {
	all := notifications.All(types.ParseReason(NotificationReason))
	printer.Notifications(all)
}
