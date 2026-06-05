package work_packages

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func Search(input string) ([]*models.WorkPackage, error) {
	filters := []requests.Filter{resources.TypeAheadFilter(input)}
	query := requests.NewFilterQuery(filters)

	response, err := requests.Get(paths.WorkPackages(), &query)
	if err != nil {
		return nil, err
	}

	collection := parser.Parse[dtos.WorkPackageCollectionDto](response)
	return collection.Convert().Items, nil
}
