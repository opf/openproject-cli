package workpackage

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "work-package [verb]",
	Short: "Manage work packages",
	Long:  "Create, list, update, inspect, and manage work packages in OpenProject.",
}

func init() {
	initListFlags()

	createCmd.Flags().Uint64VarP(
		&createProjectId,
		"project",
		"p",
		0,
		"Project ID to create the work package in",
	)
	_ = createCmd.MarkFlagRequired("project")
	createCmd.Flags().BoolVarP(
		&createOpenInBrowser,
		"open",
		"o",
		false,
		"Open the created work package in the default browser",
	)
	createCmd.Flags().StringVarP(
		&createTypeFlag,
		"type",
		"t",
		"",
		"Change the work package type",
	)

	updateCmd.Flags().StringVarP(
		&updateActionFlag,
		"action",
		"a",
		"",
		"Executes a custom action on a work package",
	)
	updateCmd.Flags().Uint64Var(
		&updateAssigneeFlag,
		"assignee",
		0,
		"Assign a user to the work package",
	)
	updateCmd.Flags().StringVar(
		&updateAttachFlag,
		"attach",
		"",
		"Attach a file to the work package",
	)
	updateCmd.Flags().StringVar(
		&updateSubjectFlag,
		"subject",
		"",
		"Change the subject of the work package",
	)
	updateCmd.Flags().StringVarP(
		&updateTypeFlag,
		"type",
		"t",
		"",
		"Change the work package type",
	)

	inspectCmd.Flags().BoolVarP(
		&inspectOpenInBrowser,
		"open",
		"o",
		false,
		"Open the work package in the default browser",
	)
	inspectCmd.Flags().BoolVar(
		&inspectListAvailableTypes,
		"types",
		false,
		"List the available types on the work package.",
	)

	RootCmd.AddCommand(listCmd, createCmd, updateCmd, inspectCmd)
}
