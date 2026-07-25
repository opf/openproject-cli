package workpackage

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "work-package [verb]",
	Short: "Manage work packages",
	Long:  "Create, list, update, inspect, and manage work packages in OpenProject.",
}

func init() {
	initListFlags()

	createCmd.Flags().StringVarP(
		&createProjectId,
		"project",
		"p",
		"",
		"Project numeric ID or identifier to create the work package in",
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
	createCmd.Flags().Uint64Var(
		&createAssigneeFlag,
		"assignee",
		0,
		"Assign a user to the work package",
	)
	createCmd.Flags().StringVar(
		&createDescriptionFlag,
		"description",
		"",
		"Description of the work package (markdown)",
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
		&updateDescriptionFlag,
		"description",
		"",
		"Description of the work package (markdown)",
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
	updateCmd.Flags().StringVar(
		&updateStatusFlag,
		"status",
		"",
		"Change the status of the work package by name",
	)
	updateCmd.Flags().StringArrayVar(
		&updateSetFlags,
		"set",
		nil,
		"Set a field by label or API name, e.g. --set \"Story points=5\"",
	)
	updateCmd.Flags().BoolVar(
		&updateDryRun,
		"dry-run",
		false,
		"Validate and show the resulting plan without applying changes",
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
	inspectCmd.Flags().BoolVar(
		&inspectWithChildren,
		"children",
		false,
		"Include direct children and full field details in the output",
	)

	searchCmd.Flags().StringVarP(
		&searchProjectId,
		"project",
		"p",
		"",
		"Limit search to a project (numeric ID or identifier)",
	)

	RootCmd.AddCommand(listCmd, createCmd, updateCmd, inspectCmd, searchCmd)
}
