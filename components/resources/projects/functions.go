package projects

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/requests"
)

const path = "api/v3/projects"

func All() []*Project {
	response := requests.Get(path)
	element := parser.Parse[ProjectCollectionDto](response)
	return element.convert()
}
