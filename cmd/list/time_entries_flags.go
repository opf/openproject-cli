package list

func initTimeEntriesFlags() {
	timeEntriesCmd.Flags().StringVarP(
		&user,
		"user",
		"u",
		"me",
		"User the time entry tracks expenditures for (can be name, ID or 'me')",
	)
}
