package time_entries

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func All(query requests.Query) ([]*models.TimeEntry, error) {
	response, err := requests.Get(paths.TimeEntries(), &query)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.TimeEntryCollectionDto](response)
	return element.Convert(), nil
}
