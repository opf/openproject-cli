package notification

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "notification [verb]",
	Short: "Manage notifications",
	Long:  "List notifications in OpenProject.",
}

func init() {
	listCmd.Flags().StringVarP(
		&listReason,
		"reason",
		"r",
		"",
		"The reason for the notification",
	)

	RootCmd.AddCommand(listCmd)
}
