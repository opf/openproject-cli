package project

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "project [verb]",
	Short: "Manage projects",
	Long:  "List and inspect projects in OpenProject.",
}

func init() {
	inspectCmd.Flags().BoolVarP(
		&openInBrowser,
		"open",
		"o",
		false,
		"Open the project in the default browser",
	)

	RootCmd.AddCommand(listCmd, inspectCmd)
}
