package version

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "version [verb]",
	Short: "Manage versions",
	Long:  "List a project's versions in OpenProject.",
}

func init() {
	listCmd.Flags().StringVarP(
		&listProjectId,
		"project",
		"p",
		"",
		"Project numeric ID or identifier",
	)

	_ = listCmd.MarkFlagRequired("project")

	RootCmd.AddCommand(listCmd)
}
