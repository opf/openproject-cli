package list

func initTimeEntriesFlags() {
	for _, filter := range activeTimeEntryFilters {
		timeEntriesCmd.Flags().StringVarP(
			filter.ValuePointer(),
			filter.Name(),
			filter.ShortHand(),
			filter.DefaultValue(),
			filter.Usage(),
		)
	}
}
