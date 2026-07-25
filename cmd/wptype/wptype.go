package wptype

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "type [verb]",
	Short: "Manage work package types",
	Long:  "List work package types in OpenProject.",
}

func init() {
	RootCmd.AddCommand(listCmd)
}
