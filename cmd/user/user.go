package user

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "user [verb]",
	Short: "Manage users",
	Long:  "Search and inspect users in OpenProject.",
}

func init() {
	RootCmd.AddCommand(searchCmd)
}
