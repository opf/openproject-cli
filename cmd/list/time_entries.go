package list

import (
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/resources/time_entries"
	"github.com/spf13/cobra"
)

var timeEntriesCmd = &cobra.Command{
	Use:   "timeentries",
	Short: "Lists time entries",
	Long:  "Get a list of all personal time entries.",
	Run:   listTimeEntries,
}

func listTimeEntries(_ *cobra.Command, _ []string) {
	if all, err := time_entries.All(); err == nil {
		printer.TimeEntryList(all)
	} else {
		printer.Error(err)
	}
}
