package list

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "list [resource]",
	Short: "Lists the specific resource",
	Long: `Get a list of the ordered resource.
The list can get filtered further.`,
}

func init() {
	RootCmd.AddCommand(projectsCmd)
}
