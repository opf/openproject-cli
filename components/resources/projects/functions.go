package projects

import (
	"path/filepath"
	"strconv"

	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

const path = "api/v3/projects"

func All() ([]*models.Project, error) {
	response, err := requests.Get(path, nil)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.ProjectCollectionDto](response)
	return element.Convert(), nil
}

func Lookup(id uint64) (*models.Project, error) {
	response, err := requests.Get(filepath.Join(path, strconv.FormatUint(id, 10)), nil)
	if err != nil {
		return nil, err
	}

	element := parser.Parse[dtos.ProjectDto](response)
	return element.Convert(), nil
}
