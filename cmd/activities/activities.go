package activities

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "activities [verb]",
	Short: "Manage activities",
	Long:  "List activities scoped by work package, project, or globally.",
}

func init() {
	listCmd.Flags().StringVarP(
		&listWpId,
		"work-package",
		"w",
		"",
		"Work package ID or identifier to list activities for",
	)
	_ = listCmd.MarkFlagRequired("work-package")

	RootCmd.AddCommand(listCmd)
}
