package work_packages

import (
	"strconv"
	"strings"

	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
)

func availableTypes(workPackage *dtos.WorkPackageDto) ([]*dtos.TypeDto, error) {
	response, err :=requests.Get(workPackage.Links.Project.Href, nil)
	if err != nil {
		return nil, err
	}

	project := parser.Parse[dtos.ProjectDto](response)
	response, err = requests.Get(project.Links.Types.Href, nil)
	if err != nil {
		return nil, err
	}

	return parser.Parse[dtos.TypeCollectionDto](response).Embedded.Elements, nil
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
