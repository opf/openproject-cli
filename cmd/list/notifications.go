package list

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/notifications"
)

var NotificationReason string

var validReasons = []string{"", "assigned", "mentioned"}
var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Lists notifications",
	Long: `Get a list of unread notifications.
The list can get filtered further.`,
	Run: listNotifications,
}

func listNotifications(_ *cobra.Command, _ []string) {
	if !common.Contains(validReasons, NotificationReason) {
		printer.ErrorText(fmt.Sprintf("Reason '%s' is invalid.", NotificationReason))
	}

	all := notifications.All(NotificationReason)
	printer.Notifications(all)
}
