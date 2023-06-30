package work_packages

import (
	"strings"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
)

func availableTypes(projectLink *dtos.LinkDto) (dtos.TypeDtos, error) {
	response, err := requests.Get(projectLink.Href, nil)
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

func findType(input string, availableTypes dtos.TypeDtos) *dtos.TypeDto {
	typeAsId, typeId := common.ParseId(input)

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
