package time_entries

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func All() ([]*models.TimeEntry, error) {
	var filters []requests.Filter
	filters = append(filters, UserFilter("me"))

	query := requests.NewPaginatedQuery(-1, filters)

	response, err := requests.Get(paths.TimeEntries(), &query)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.TimeEntryCollectionDto](response)
	return element.Convert(), nil
}
