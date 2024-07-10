package start

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "start [resource]",
	Short: "Starts the work on a resource",
	Long:  "Creates a branch based on the given resource and switches to it.",
}

func init() {
	RootCmd.AddCommand(startWorkPackageCmd)
}
