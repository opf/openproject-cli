package budget

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "budget [verb]",
	Short: "Manage budgets",
	Long:  "List and inspect budgets in OpenProject.",
}

func init() {
	listCmd.Flags().Uint64VarP(
		&listProjectId,
		"project",
		"p",
		0,
		"Project id",
	)

	_ = listCmd.MarkFlagRequired("project")

	RootCmd.AddCommand(listCmd, inspectCmd)
}
