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
	workPackageCmd.Flags().StringVar(
		&descriptionFlag,
		"description",
		"",
		"Change the raw work package description",
	)
	workPackageCmd.Flags().StringVar(
		&statusFlag,
		"status",
		"",
		"Change the status of the work package (by name)",
	)
	workPackageCmd.Flags().StringVarP(
		&typeFlag,
		"type",
		"t",
		"",
		"Change the work package type",
	)
	workPackageCmd.Flags().StringArrayVar(
		&setFlags,
		"set",
		nil,
		"Set a schema-resolved custom field using label=value or apiField=value",
	)
	workPackageCmd.Flags().BoolVar(
		&printUpdatedWorkPackageAsJSON,
		"json",
		false,
		"Print machine-readable JSON output",
	)
	workPackageCmd.Flags().BoolVar(
		&dryRunUpdateWorkPackage,
		"dry-run",
		false,
		"Resolve and validate without persisting the update",
	)
}
