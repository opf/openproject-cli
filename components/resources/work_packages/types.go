package work_packages

import (
	"strconv"
	"strings"

	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
)

func availableTypes(workPackage *dtos.WorkPackageDto) []*dtos.TypeDto {
	status, response := requests.Get(workPackage.Links.Project.Href, nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	project := parser.Parse[dtos.ProjectDto](response)
	status, response = requests.Get(project.Links.Types.Href, nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	return parser.Parse[dtos.TypeCollectionDto](response).Embedded.Elements
}

func findType(input string, availableTypes []*dtos.TypeDto) *dtos.TypeDto {
	var typeAsId = false
	typeId, err := strconv.ParseUint(input, 10, 64)
	if err == nil {
		typeAsId = true
	}

	var found []*dtos.TypeDto
	for _, t := range availableTypes {
		if typeAsId && typeId == uint64(t.Id) ||
			!typeAsId && strings.ToLower(input) == strings.ToLower(t.Name) {
			found = append(found, t)
		}
	}

	if len(found) == 1 {
		return found[0]
	}

	return nil
}
