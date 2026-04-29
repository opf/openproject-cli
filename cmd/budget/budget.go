package budget

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "budget [verb]",
	Short: "Manage budgets",
	Long:  "List and inspect budgets in OpenProject.",
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

	RootCmd.AddCommand(listCmd, inspectCmd)
}
