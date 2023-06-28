package projects

import (
	"path/filepath"
	"strconv"

	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func AvailableVersions(projectId uint64) ([]*models.Version, error) {
	response, err := requests.Get(filepath.Join(path, strconv.FormatUint(projectId, 10), "versions"), nil)
	if err != nil {
		return nil, err
	}

	versionCollection := parser.Parse[dtos.VersionCollectionDto](response)

	return versionCollection.Convert(), nil
}
