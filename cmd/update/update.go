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
	addWorkPackageFlags()

	RootCmd.AddCommand(workPackageCmd)
}

func addWorkPackageFlags() {
	workPackageCmd.Flags().StringVarP(
		&actionFlag,
		"action",
		"a",
		"",
		"Executes a custom action on a work package",
	)
	workPackageCmd.Flags().Uint64Var(
		&assigneeFlag,
		"assignee",
		0,
		"Assign a user to the work package",
	)
	workPackageCmd.Flags().StringVar(
		&attachFlag,
		"attach",
		"",
		"Attach a file to the work package",
	)
	workPackageCmd.Flags().StringVar(
		&subjectFlag,
		"subject",
		"",
		"Change the subject of the work package",
	)
	workPackageCmd.Flags().StringVarP(
		&typeFlag,
		"type",
		"t",
		"",
		"Change the work package type",
	)
}
