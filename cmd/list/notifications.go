package list

import (
	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/notifications"
)

var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Lists notifications",
	Long: `Get a list of unread notifications.
The list can get filtered further.`,
	Run: listNotifications,
}

func listNotifications(cmd *cobra.Command, args []string) {
	all := notifications.All()
	printer.Notifications(all)
}
