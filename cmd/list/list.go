package list

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "list [resource]",
	Short: "Lists the specific resource",
	Long: `Get a list of the ordered resource.
The list can get filtered further.`,
}

func init() {
	notificationsCmd.Flags().StringVarP(
		&notificationReason,
		"reason",
		"r",
		"",
		"The reason for the notification",
	)

	RootCmd.AddCommand(
		projectsCmd,
		notificationsCmd,
	)
}
