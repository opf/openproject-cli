package projects

import (
	"path/filepath"
	"strconv"

	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/models"
)

const path = "api/v3/projects"

func All() []*models.Project {
	status, response := requests.Get(path, nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	element := parser.Parse[ProjectCollectionDto](response)
	return element.convert()
}

func Find(id int) *models.Project {
	status, response := requests.Get(filepath.Join(path, strconv.Itoa(id)), nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	element := parser.Parse[ProjectDto](response)
	return element.convert()
}
