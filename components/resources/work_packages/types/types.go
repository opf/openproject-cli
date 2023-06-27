package work_package_types

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func AvailableTypes(workPackage *dtos.WorkPackageDto) []*models.Type {
	status, response := requests.Get(workPackage.Links.Project.Href, nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}
	
	project := parser.Parse[dtos.ProjectDto](response)
	status, response = requests.Get(project.Links.Types.Href, nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	t := parser.Parse[dtos.TypeCollectionDto](response)
	return t.Convert()
}