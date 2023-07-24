package inspect

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "inspect [type] [id]",
	Short: "Show details about an object",
	Long:  "Show detailed information of an object of a specific type referenced by it's ID.",
}

func init() {
	inspectProjectCmd.Flags().BoolVarP(
		&shouldOpenProjectInBrowser,
		"open",
		"o",
		false,
		"Open the project in the default browser",
	)

	inspectWorkPackageCmd.Flags().BoolVarP(
		&shouldOpenWorkPackageInBrowser,
		"open",
		"o",
		false,
		"Open the work package in the default browser",
	)

	inspectWorkPackageCmd.Flags().BoolVar(
		&listAvailableTypes,
		"types",
		false,
		"List the available types on the work package.",
	)

	RootCmd.AddCommand(inspectProjectCmd, inspectWorkPackageCmd)
}
