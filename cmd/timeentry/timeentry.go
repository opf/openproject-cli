package timeentry

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "time-entry [verb]",
	Short: "Manage time entries",
	Long:  "List time entries in OpenProject.",
}

func init() {
	initListFlags()

	RootCmd.AddCommand(listCmd)
}
