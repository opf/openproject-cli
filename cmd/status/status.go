package status

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "status [verb]",
	Short: "Manage statuses",
	Long:  "List work package statuses in OpenProject.",
}

func init() {
	RootCmd.AddCommand(listCmd)
}
