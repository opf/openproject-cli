package work_packages

import (
	"path/filepath"
	"strconv"

	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func Activities(id uint64) (activites []*models.Activity, err error) {
	status, response := requests.Get(filepath.Join(workPackagesPath, strconv.FormatUint(id, 10), "activities"), nil)
	if !requests.IsSuccess(status) {
		printer.ResponseError(status, response)
	}

	activitiesDto := parser.Parse[dtos.ActivityCollectionDto](response)
	return activitiesDto.Convert()
}
