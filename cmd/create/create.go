package create

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "create [resource]",
	Short: "Creates a specific resource",
	Long:  "Create a specific resource in OpenProject",
}

func init() {
	createWorkPackageCmd.Flags().Uint64VarP(
		&projectId,
		"project",
		"p",
		0,
		"Project ID to create the work package in",
	)

	createWorkPackageCmd.Flags().BoolVarP(
		&shouldOpenWorkPackageInBrowser,
		"open",
		"o",
		false,
		"Open the created work package in the default browser",
	)

	createWorkPackageCmd.Flags().StringVarP(
		&typeFlag,
		"type",
		"t",
		"",
		"Change the work package type",
	)

	createWorkPackageCmd.Flags().Uint64Var(
		&parentWorkPackageID,
		"parent",
		0,
		"Create the work package as a child of an existing work package",
	)

	createWorkPackageCmd.Flags().BoolVar(
		&printCreatedWorkPackageAsJSON,
		"json",
		false,
		"Print machine-readable JSON output",
	)

	createWorkPackageCmd.Flags().BoolVar(
		&dryRunCreateWorkPackage,
		"dry-run",
		false,
		"Resolve and validate without persisting the work package",
	)

	RootCmd.AddCommand(createWorkPackageCmd)
}
