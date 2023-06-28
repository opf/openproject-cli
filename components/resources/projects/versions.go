package projects

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func AvailableVersions(projectId uint64) ([]*models.Version, error) {
	response, err := requests.Get(paths.ProjectVersions(projectId), nil)
	if err != nil {
		return nil, err
	}

	versionCollection := parser.Parse[dtos.VersionCollectionDto](response)

	return versionCollection.Convert(), nil
}
