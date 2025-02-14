package list

import (
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources"
	"github.com/opf/openproject-cli/components/resources/time_entries"
	"github.com/opf/openproject-cli/components/resources/time_entries/filters"
	"github.com/spf13/cobra"
)

var activeTimeEntryFilters = map[string]resources.Filter{
	"user": filters.NewUserFilter(),
}

var timeEntriesCmd = &cobra.Command{
	Use:   "timeentries",
	Short: "Lists time entries",
	Long:  "Get a list of all time entries.",
	Run:   listTimeEntries,
}

func listTimeEntries(_ *cobra.Command, _ []string) {
	query, err := buildTimeEntriesQuery()
	if err != nil {
		printer.ErrorText(err.Error())
		return
	}

	if all, err := time_entries.All(query); err == nil {
		printer.TimeEntryList(all)
	} else {
		printer.Error(err)
	}
}

func buildTimeEntriesQuery() (requests.Query, error) {
	var q requests.Query

	for _, filter := range activeTimeEntryFilters {
		err := filter.ValidateInput()
		if err != nil {
			return requests.NewEmptyQuery(), err
		}

		q = q.Merge(filter.Query())
	}

	return q, nil
}
