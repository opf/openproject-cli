package update

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "update [resource] [id]",
	Short: "Updates the specific resource",
	Long: `Sends an update to the given resource,
which is identified by its id. The data
to update is determined by the provided
flags.`,
}

func init() {
	workPackageCmd.Flags().StringVarP(
		&actionFlag,
		"action",
		"a",
		"",
		"Executes a custom action on a work package",
	)
	workPackageCmd.Flags().StringVar(
		&attachFlag,
		"attach",
		"",
		"Attach a file to the work package",
	)

	RootCmd.AddCommand(workPackageCmd)
}
