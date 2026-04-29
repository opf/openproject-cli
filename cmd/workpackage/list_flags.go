package workpackage

func initListFlags() {
	listCmd.Flags().StringVarP(
		&listAssignee,
		"assignee",
		"a",
		"",
		"Assignee of the work package (can be ID or 'me')",
	)

	listCmd.Flags().Uint64VarP(
		&listParentId,
		"parent-id",
		"",
		0,
		"Show only direct children of the specified work package ID")

	listCmd.Flags().StringVarP(
		&listProjectId,
		"project",
		"p",
		"",
		"Show only work packages within the specified project (numeric ID or identifier)")

	listCmd.Flags().StringVarP(
		&listStatusFilter,
		"status",
		"s",
		"",
		`Show only work packages having the specified status. The value can be the
keywords 'open', 'closed', a single ID or a comma separated array of IDs, i.e.
'7,13'. Multiple values are concatenated with a logical 'OR'. If the IDs are
prefixed with an '!' the list is instead filtered to not have the specified
status.`)

	listCmd.Flags().StringVarP(
		&listTypeFilter,
		"type",
		"t",
		"",
		`Show only work packages having the specified types. The value can be a single
ID or a comma separated array of IDs, i.e. '7,13'. Multiple values are
concatenated with a logical 'OR'. If the IDs are prefixed with an '!' the list
is instead filtered to not have the specified status.`)

	listCmd.Flags().BoolVarP(
		&listIncludeSubProjects,
		"include-sub-projects",
		"",
		false,
		`If listing the work packages of a project, this flag indicates if work
packages of sub projects should be included in the list. If omitting the flag,
the default is false.`)

	listCmd.Flags().BoolVarP(
		&listShowTotal,
		"total",
		"",
		false,
		"Show only the total number of work packages matching the filter options.")

	for _, filter := range activeFilters {
		listCmd.Flags().StringVarP(
			filter.ValuePointer(),
			filter.Name(),
			filter.ShortHand(),
			filter.DefaultValue(),
			filter.Usage(),
		)
	}
}
