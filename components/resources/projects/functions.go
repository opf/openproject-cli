package projects

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func All() ([]*models.Project, error) {
	query := requests.NewPagedQuery(-1, nil)
	response, err := requests.Get(paths.Projects(), &query)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.ProjectCollectionDto](response)
	return element.Convert(), nil
}

func Lookup(id uint64) (*models.Project, error) {
	response, err := requests.Get(paths.Project(id), nil)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.ProjectDto](response)
	return element.Convert(), nil
}
