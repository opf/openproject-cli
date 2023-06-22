package update

import "github.com/spf13/cobra"

var actionFlag string

var workPackageCmd = &cobra.Command{
	Use:   "workpackage [id]",
	Short: "Updates the work package",
	Long: `Get a list of unread notifications.
The list can get filtered further.`,
	Run: updateWorkPackage,
}

func updateWorkPackage(cmd *cobra.Command, args []string)  {
	
}
