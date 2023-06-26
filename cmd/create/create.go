package create

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "create [resource]",
	Short: "Creates a specific resource",
	Long:  "Create a specific resource in OpenProject",
}

func init() {
	RootCmd.AddCommand(createWorkPackageCmd)
}
