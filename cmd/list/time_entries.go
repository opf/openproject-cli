package list

import (
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/time_entries"
	"github.com/spf13/cobra"
)

var user string

var timeEntriesCmd = &cobra.Command{
	Use:   "timeentries",
	Short: "Lists time entries",
	Long:  "Get a list of all time entries.",
	Run:   listTimeEntries,
}

func listTimeEntries(_ *cobra.Command, _ []string) {
	if all, err := time_entries.All(timeEntriesFilterOptions()); err == nil {
		printer.TimeEntryList(all)
	} else {
		printer.Error(err)
	}
}

func timeEntriesFilterOptions() *map[time_entries.FilterOption]string {
	options := make(map[time_entries.FilterOption]string)

	if len(user) > 0 {
		options[time_entries.User] = user
	}

	return &options
}
