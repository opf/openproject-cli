package list

func initWorkPackagesFlags() {
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
		&statusFilter,
		"status",
		"s",
		"",
		`Show only work packages having the specified status. The value can be the
keywords 'open', 'closed', a single ID or a comma separated array of IDs, i.e.
'7,13'. Multiple values are concatenated with a logical 'OR'. If the IDs are
prefixed with an '!' the list is instead filtered to not have the specified
status.`)

	workPackagesCmd.Flags().StringVarP(
		&typeFilter,
		"type",
		"t",
		"",
		`Show only work packages having the specified types. The value can be a single
ID or a comma separated array of IDs, i.e. '7,13'. Multiple values are
concatenated with a logical 'OR'. If the IDs are prefixed with an '!' the list
is instead filtered to not have the specified status.`)

	workPackagesCmd.Flags().BoolVarP(
		&includeSubProjects,
		"include-sub-projects",
		"",
		false,
		`If listing the work packages of a project, this flag indicates if work
packages of sub projects should be included in the list. If omitting the flag,
the default is false.`)

	workPackagesCmd.Flags().BoolVarP(
		&showTotal,
		"total",
		"",
		false,
		"Show only the total number of work packages matching the filter options.")

	for _, filter := range activeFilters {
		workPackagesCmd.Flags().StringVarP(
			filter.ValuePointer(),
			filter.Name(),
			filter.ShortHand(),
			filter.DefaultValue(),
			filter.Usage(),
		)
	}
}
