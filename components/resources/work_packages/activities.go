package work_packages

import (
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/dtos"
	"github.com/opf/openproject-cli/models"
)

func Activities(id uint64) (activites []*models.Activity, err error) {
	response, err := requests.Get(paths.WorkPackageActivities(id), nil)
	if err != nil {
		printer.Error(err)
	}

	activitiesDto := parser.Parse[dtos.ActivityCollectionDto](response)
	return activitiesDto.Convert()
}
