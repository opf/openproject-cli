package create

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "create [resource]",
	Short: "Creates a specific resource",
	Long:  "Create a specific resource in OpenProject",
}

func init() {
	createWorkPackageCmd.Flags().Int64VarP(
		&projectId,
		"project",
		"p",
		0,
		"Project ID to create the work package in",
	)
	_ = createWorkPackageCmd.MarkFlagRequired("project")

	createWorkPackageCmd.Flags().BoolVarP(
		&shouldOpenWorkPackageInBrowser,
		"open",
		"o",
		false,
		"Open the created work package in the default browser",
	)

	RootCmd.AddCommand(createWorkPackageCmd)
}
