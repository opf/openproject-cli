package timeentry

func initListFlags() {
	for _, filter := range activeTimeEntryFilters {
		listCmd.Flags().StringVarP(
			filter.ValuePointer(),
			filter.Name(),
			filter.ShortHand(),
			filter.DefaultValue(),
			filter.Usage(),
		)
	}
}
