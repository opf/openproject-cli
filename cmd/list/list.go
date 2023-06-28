package list

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "list [resource]",
	Short: "Lists the specific resource",
	Long: `Get a list of the ordered resource.
The list can get filtered further.`,
}

func init() {
	notificationsCmd.Flags().StringVarP(
		&notificationReason,
		"reason",
		"r",
		"",
		"The reason for the notification",
	)

	workPackagesCmd.Flags().StringVarP(
		&assignee,
		"assignee",
		"a",
		"",
		"Assignee of the work package (can be name, ID or 'me')",
	)

	workPackagesCmd.Flags().Uint64VarP(
		&projectId,
		"project-id",
		"p",
		0,
		"Show only work packages within the specified projectId")

	workPackagesCmd.Flags().StringVarP(
		&version,
		"version",
		"v",
		"",
		"Show only work packages having the specified version")

	RootCmd.AddCommand(
		projectsCmd,
		notificationsCmd,
		workPackagesCmd,
		activitiesCmd,
	)
}
