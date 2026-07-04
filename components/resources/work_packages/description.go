package work_packages

import (
	"os"
	"strconv"

	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/dtos"
)

func description(input string) *dtos.LongTextDto {
	return &dtos.LongTextDto{Format: "markdown", Raw: input}
}

func descriptionFromFile(path string) (*dtos.LongTextDto, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return description(string(data)), nil
}

func parentLink(input string) (*dtos.LinkDto, error) {
	parentId, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return nil, err
	}

	return &dtos.LinkDto{Href: paths.WorkPackage(parentId)}, nil
}
