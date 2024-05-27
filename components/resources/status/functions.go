package status

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func All() ([]*models.Status, error) {
	query := requests.NewPaginatedQuery(-1, nil)
	response, err := requests.Get(paths.Status(), &query)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.StatusCollectionDto](response)
	return element.Convert(), nil
}
