package timeentry

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "time-entry [verb]",
	Short: "Manage time entries",
	Long:  "List and create time entries in OpenProject.",
}

func init() {
	initListFlags()
	initCreateFlags()

	RootCmd.AddCommand(listCmd, createCmd)
}
